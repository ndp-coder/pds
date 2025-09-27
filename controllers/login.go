package controllers

import (
	"PDS/database"
	"PDS/midlewares"
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var Login_details struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func HashPassword(password string) (string, error) {
	hashpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashpass), nil
}

func VerifyHashPassword(hashedPassword, Password string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(Password))
	return err == nil
}

// Login godoc
// @Summary Login user
// @Description Authenticate a user (doctor, pharmacist, or admin)
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body models.User true "User login credentials"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func Login(c *gin.Context) {
	 // logic

	if err := c.ShouldBindJSON(&Login_details); err != nil {
		c.JSON(404, gin.H{
			"error": "Invalid input",
		})
		c.Abort()
		return
	}

	var HashedPass string

	err := database.Postdb.QueryRow(context.Background(), "select password from users where email = $1", Login_details.Email).Scan(&HashedPass)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(404, gin.H{"error": "Email not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database error", "details": err.Error()})
		}
		return
	}
	token, err := midlewares.GenerateToken(Login_details.Email)
	if err != nil {
		c.JSON(404, gin.H{
			"error": err,
		})
	}
	is_pass_correct := VerifyHashPassword(HashedPass, Login_details.Password)
	if is_pass_correct {
		c.JSON(200, gin.H{
			"message": "user logined successfully",
			"token":   token,
		})
		return
	} else {
		c.JSON(404, gin.H{
			"error": "password is incorrect",
			"pass":  HashedPass,
		})
		c.Abort()
		return
	}

}
