package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"PDS/controllers"
	"PDS/database"
	"PDS/midlewares"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Setup router with auth-protected endpoints
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	protected := r.Group("/")
	protected.Use(midlewares.AuthMiddleware())
	{
		protected.POST("/addmedicine", controllers.Add_Medicine)
		protected.POST("/Dec-med", controllers.Decrease_medicn)
	}
	return r
}

// ---------- TEST 1: Race on Add Medicine ----------
func TestRaceConditionOnAddMedicine(t *testing.T) {
	database.ConnectDB()
	ctx := context.Background()

	// Short unique medicine name
	uniqueMed := fmt.Sprintf("MedAdd%d", time.Now().UnixNano()%1000000)

	// Ensure fresh row
	_, err := database.Postdb.Exec(ctx,
		`INSERT INTO medicine (medicine_name, dosage_form, stock_quantity)
		 VALUES ($1, $2, $3)`,
		uniqueMed, "tablet", 0)
	assert.NoError(t, err)

	// Mock admin token
	token, _ := midlewares.GenerateToken("admin@pds.com")

	r := setupRouter()

	var wg sync.WaitGroup
	concurrency := 10
	addQty := 5

	// Run 10 concurrent Add_Medicine requests
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			body, _ := json.Marshal(map[string]interface{}{
				"medicine_name":  uniqueMed,
				"dosage_form":    "tablet",
				"stock_quantity": addQty,
			})

			req, _ := http.NewRequest("POST", "/addmedicine", bytes.NewBuffer(body))
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Contains(t, []int{http.StatusOK, http.StatusInternalServerError}, w.Code)
		}()
	}
	wg.Wait()

	// Check final stock (must equal concurrency * addQty)
	var finalStock int
	err = database.Postdb.QueryRow(ctx,
		`SELECT stock_quantity FROM medicine WHERE medicine_name=$1 AND dosage_form=$2`,
		uniqueMed, "tablet").Scan(&finalStock)
	assert.NoError(t, err)

	expected := concurrency * addQty
	assert.Equal(t, expected, finalStock, "Race condition detected in Add_Medicine")
}

// ---------- TEST 2: Race on Decrease Medicine ----------
func TestRaceConditionOnDecreaseMedicine(t *testing.T) {
	database.ConnectDB()
	ctx := context.Background()

	// Short unique medicine name
	uniqueMed := fmt.Sprintf("MedDec%d", time.Now().UnixNano()%1000000)

	// Insert initial stock = 50
	_, err := database.Postdb.Exec(ctx,
		`INSERT INTO medicine (medicine_name, dosage_form, stock_quantity)
		 VALUES ($1, $2, $3)`,
		uniqueMed, "capsule", 50)
	assert.NoError(t, err)

	// Mock pharmacist token
	token, _ := midlewares.GenerateToken("pharma@pds.com")

	r := setupRouter()

	var wg sync.WaitGroup
	concurrency := 10
	decreaseQty := 5

	// Run 10 concurrent Decrease_medicn requests
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			body, _ := json.Marshal(map[string]interface{}{
				"medicine_name":  uniqueMed,
				"dosage_form":    "capsule",
				"stock_quantity": decreaseQty,
			})

			req, _ := http.NewRequest("POST", "/Dec-med", bytes.NewBuffer(body))
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// OK if stock was available, BadRequest when insufficient
			assert.Contains(t, []int{http.StatusOK, http.StatusBadRequest}, w.Code)
		}()
	}
	wg.Wait()

	// Final stock must never go negative
	var finalStock int
	err = database.Postdb.QueryRow(ctx,
		`SELECT stock_quantity FROM medicine WHERE medicine_name=$1 AND dosage_form=$2`,
		uniqueMed, "capsule").Scan(&finalStock)
	assert.NoError(t, err)

	assert.True(t, finalStock >= 0, "Stock went negative, race condition!")
}
