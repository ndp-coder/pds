	package main

	import (
		"PDS/controllers"
		"PDS/database"
		"PDS/midlewares"
		"time"

		_ "PDS/docs" // Swagger docs

		"github.com/gin-contrib/cors"
		"github.com/gin-gonic/gin"
		swaggerFiles "github.com/swaggo/files"
		ginSwagger "github.com/swaggo/gin-swagger"
	)

	// @title Prescription Dispensing System API
	// @version 1.0
	// @description This is the backend API for managing doctors, pharmacists, medicines, and prescriptions.
	// @termsOfService http://swagger.io/terms/

	// @contact.name NDP
	// @contact.email support@pds.com
	// @securityDefinitions.apikey ApiKeyAuth
	// @in header
	// @name Authorization


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

		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"}, // Allow all origins, or specify your frontend URL
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))

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
