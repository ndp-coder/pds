package midlewares

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func getJWTKey() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func GenerateToken(email string) (string, error) {
	expirationTime := time.Now().Add(128 * time.Hour)
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTKey())
}



func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHead := c.GetHeader("Authorization")
		if authHead == "" {
			c.JSON(401, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHead, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid Authorization format"})
			c.Abort()
			return
		}

		realToken := strings.TrimSpace(parts[1])
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(realToken, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return getJWTKey(), nil // ✅ dynamic key
		})

		if err != nil {
			fmt.Println("JWT parsing error:", err)
			c.JSON(401, gin.H{
				"error":   "Invalid token",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		if !token.Valid {
			fmt.Println("JWT is invalid")
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Next()
	}
}
