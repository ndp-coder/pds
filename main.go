package main

import (
	"PDS/controllers"
	"PDS/database"
	"PDS/midlewares"

	"github.com/gin-gonic/gin"
)

func main() {

	database.ConnectDB()

	router := gin.Default()
	router.GET("/", func (c *gin.Context){
		c.JSON(200 , gin.H{
			"message" : "pds is running",
		})
	})

	router.POST("/login", controllers.Login)

	Protected := router.Group("/")
	Protected.Use(midlewares.AuthMiddleware())
	{

		Protected.POST("/createDocter", controllers.CreateDocter)
		Protected.POST("/createPharmasist", controllers.CreatePharmasist)
		Protected.POST("/addmedicine", controllers.Add_Medicine)
		Protected.POST("/prescription", controllers.Prescription)
		Protected.POST("/Dec-med", controllers.Decrease_medicn)
		Protected.GET("/getmedicen", controllers.GetMedicen)
	}

	router.Run(":8080")

}
