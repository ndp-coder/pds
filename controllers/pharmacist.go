package controllers

import (
	"PDS/database"
	"PDS/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Decrease_medicn(c *gin.Context) {
	var Prescription struct {
		Medicine_Name  string `json:"medicine_name"`
		Dosage_form    string `json:"dosage_form"`
		Stock_Quantity int    `json:"stock_quantity"`
	}

	if err := c.ShouldBindJSON(&Prescription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	emailVal, exists := c.Get("email")
    if !exists {
        c.JSON(401, gin.H{
			"error": "Email not found in context",
			"email" : emailVal,
	})
        return
    }

    email := emailVal.(string)
    c.JSON(200, gin.H{
        "message": "Welcome!",
        "email":   email,
    })

	var IsAdmin_db int 
	res := database.Postdb.QueryRow(context.Background(), "SELECT CASE WHEN EXISTS ( SELECT 1 FROM users WHERE email = $1 AND role = 'pharmasist') THEN 1 ELSE 0 END", email).Scan(&IsAdmin_db)

	if res != nil{
		c.JSON(400 , gin.H{
			"errro" : "cant fectch admin from database",
		})
		c.Abort()
		return
	}

	if IsAdmin_db == 0{
		c.JSON(404 , gin.H{
			"message" : "you are not allowed to add users",
			"email" : email,
			
		})
		c.Abort()
		return
	}


	if Prescription.Stock_Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Stock quantity must be greater than zero",
		})
		return
	}

	ctx := context.Background()

	tx, err := database.Postdb.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to start transaction",
			"details": err.Error(),
		})
		return
	}
	defer func() {

		_ = tx.Rollback(ctx)
	}()

	var currentStock int

	err = tx.QueryRow(ctx, `SELECT stock_quantity FROM medicine WHERE medicine_name = $1 AND dosage_form = $2 FOR UPDATE`,
		Prescription.Medicine_Name, Prescription.Dosage_form).Scan(&currentStock)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Medicine not found",
			"details": err.Error(),
		})
		return
	}

	if currentStock < Prescription.Stock_Quantity {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":           "Insufficient stock",
			"available_stock": currentStock,
		})
		return
	}

	var newStock int
	err = tx.QueryRow(ctx, ` UPDATE medicine SET stock_quantity = stock_quantity - $1 WHERE medicine_name = $2 AND dosage_form = $3 RETURNING stock_quantity`,
		Prescription.Stock_Quantity, Prescription.Medicine_Name, Prescription.Dosage_form).Scan(&newStock)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update stock",
			"details": err.Error(),
		})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to commit transaction",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Stock reduced successfully",
		"new_stock":   newStock,
		"medicine":    Prescription.Medicine_Name,
		"dosage_form": Prescription.Dosage_form,
	})
}

func GetMedicen(c *gin.Context) {
	res, err := database.Postdb.Query(context.Background(), "select medicine_name , dosage_form , stock_quantity from medicine")
	if err != nil {
		c.JSON(404, gin.H{"erro": err})
	}
	defer res.Close()
	var data []models.Add_Medicine
	for res.Next() {
		var a models.Add_Medicine
		res.Scan(&a.Medicine_Name, &a.Dosage_form, &a.Stock_Quantity)
		data = append(data, a)
	}
	c.JSON(200, gin.H{"data": data})
}
