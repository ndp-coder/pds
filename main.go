package main

import (
	"PDS/controllers"
	"PDS/database"
	"PDS/midlewares"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "PDS/docs" // Swagger docs
)

// @title Prescription Dispensing System API
// @version 1.0
// @description This is the backend API for managing doctors, pharmacists, medicines, and prescriptions.
// @termsOfService http://swagger.io/terms/

// @contact.name NDP
// @contact.email support@pds.com

// @host localhost:8080
// @BasePath /
func main() {
	database.ConnectDB()

	router := gin.Default()

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Root health check
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "PDS is running",
		})
	})

	// Public endpoint
	router.POST("/login", controllers.Login)

	// Protected endpoints
	protected := router.Group("/")
	protected.Use(midlewares.AuthMiddleware())
	{
		protected.POST("/createDoctor", controllers.CreateDocter)
		protected.POST("/createPharmacist", controllers.CreatePharmasist)
		protected.POST("/addmedicine", controllers.Add_Medicine)
		protected.POST("/prescription", controllers.Prescription)
		protected.POST("/Dec-med", controllers.Decrease_medicn)
		protected.GET("/getmedicen", controllers.GetMedicen)
	}

	router.Run(":8080")
}

