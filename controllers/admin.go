package controllers

import (
	"PDS/database"
	"context"

	"net/http"

	"github.com/gin-gonic/gin"
)

// AddMedicine godoc
// @Summary Add new medicine
// @Description Add a new medicine to the stock
// @Tags Medicine
// @Accept json
// @Produce json
// @Param medicine body models.Add_Medicine true "Medicine Data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Security ApiKeyAuth
// @Router /addmedicine [post]
func Add_Medicine(c *gin.Context) {
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
			"error":   "failed to start transaction",
			"err": err.Error(),
		})
		return
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var newStock int
	err = tx.QueryRow(ctx,
		`INSERT INTO medicine (medicine_name, dosage_form, stock_quantity)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (medicine_name, dosage_form)
		 DO UPDATE SET stock_quantity = medicine.stock_quantity + EXCLUDED.stock_quantity
		 RETURNING stock_quantity`,
		Prescription.Medicine_Name, Prescription.Dosage_form, Prescription.Stock_Quantity).Scan(&newStock)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to insert/update stock",
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
		"message":     "Stock added successfully",
		"new_stock":   newStock,
		"medicine":    Prescription.Medicine_Name,
		"dosage_form": Prescription.Dosage_form,
	})
}
