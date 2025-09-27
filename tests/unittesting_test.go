package tests

import (
	"PDS/database"
	"PDS/midlewares"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// Setup test DB
func init() {
	// Load .env
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	// Connect DB
	dbURL := os.Getenv("DB_DSN")
	if dbURL == "" {
		panic("DB_DSN not found in .env")
	}

	database.ConnectDB()
}

// Seed test user
func seedTestUser(t *testing.T, email, role string) {
	// Delete any existing record to avoid conflicts
	_, _ = database.Postdb.Exec(context.Background(), `DELETE FROM users WHERE email = $1`, email)

	// Insert fresh user with all required fields
	_, err := database.Postdb.Exec(
		context.Background(),
		`INSERT INTO users (name, email, role, password) VALUES ($1, $2, $3, $4)`,
		"Test User",           // name
		email,                 // email
		role,                  // role
		"hashed_password_123", // password
	)

	if err != nil {
		t.Fatalf("Failed to seed test user: %v", err)
	}
}



// Dummy handler for pharmacist check
func pharmacistOnlyHandler(c *gin.Context) {
	email := c.GetString("email")

	var isPharmacist int
	err := database.Postdb.QueryRow(context.Background(),
		`SELECT CASE 
			WHEN EXISTS (SELECT 1 FROM users WHERE email = $1 AND role = 'Pharmasist') 
			THEN 1 ELSE 0 END`,
		email,
	).Scan(&isPharmacist)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	if isPharmacist == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to add users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pharmacist authorized"})
}

// Test case
func TestPharmacistAuthorization(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(midlewares.AuthMiddleware()) // Add your JWT middleware
	router.GET("/test-pharmacist", pharmacistOnlyHandler)

	testEmail := "pharma@example.com"
	seedTestUser(t, testEmail, "Pharmasist")

	// Generate JWT token
	token, err := midlewares.GenerateToken(testEmail)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Make HTTP request
	req, _ := http.NewRequest("GET", "/test-pharmacist", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Validate response
	assert.Equal(t, http.StatusOK, rr.Code, "Expected HTTP 200 OK")
	assert.Contains(t, rr.Body.String(), "Pharmacist authorized")

	// Now test unauthorized user
	seedTestUser(t, "norr@example.com", "Customer")
	token2, _ := midlewares.GenerateToken("noer@example.com")

	req2, _ := http.NewRequest("GET", "/test-pharmacist", nil)
	req2.Header.Set("Authorization", "Bearer "+token2)

	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	assert.Equal(t, http.StatusForbidden, rr2.Code, "Expected HTTP 403 Forbidden")
	assert.Contains(t, rr2.Body.String(), "not allowed")
}
