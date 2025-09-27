package controllers

import (
	"PDS/database"
	"context"

	"github.com/gin-gonic/gin"
)

// Prescription godoc
// @Summary Create prescription
// @Description Doctor can create a prescription for a patient
// @Tags Prescription
// @Accept json
// @Produce json
// @Param prescription body models.Prescription true "Prescription Data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Security ApiKeyAuth
// @Router /prescription [post]
func Prescription(c *gin.Context) {
	 // logic

	var Pres_details struct {
		Patient_Name  string `json:"patient_name"`
		Medicine_name string `json:"medicine_name"`
		Dosage_form   string `json:"dosage_form"`
		Quantity      int    `json:"quantity"`
		Ph_num        string `json:"ph_num"`
	}

	if err := c.ShouldBindJSON(&Pres_details); err != nil {
		c.JSON(404, gin.H{
			"error": "Invalid input",
		})
	}

	emailVal, exists := c.Get("email")
	if !exists {
		c.JSON(401, gin.H{
			"error": "Email not found in context",
			"email": emailVal,
		})
		return
	}

	email := emailVal.(string)
	c.JSON(200, gin.H{
		"message": "Welcome!",
		"email":   email,
	})

	var IsAdmin_db int
	res_0 := database.Postdb.QueryRow(context.Background(), "SELECT CASE WHEN EXISTS ( SELECT 1 FROM users WHERE email = $1 AND role = 'docter') THEN 1 ELSE 0 END", email).Scan(&IsAdmin_db)

	if res_0 != nil {
		c.JSON(400, gin.H{
			"errro": "cant fectch admin from database",
		})
		c.Abort()
		return
	}

	if IsAdmin_db == 0 {
		c.JSON(404, gin.H{
			"message": "you are not allowed to add users",
			"email":   email,
		})
		c.Abort()
		return
	}

	var is_user_in_prescriptions int

	database.Postdb.QueryRow(context.Background(), "select ph_num from prescriptions where ph_num = ? and patient_name = ?", Pres_details.Ph_num, Pres_details.Patient_Name).Scan(&is_user_in_prescriptions)

	

	if is_user_in_prescriptions != 0 {
		database.Postdb.Exec(context.Background(), "delete from prescriptions where ph_num = ? and patient_name = ?", Pres_details.Ph_num, Pres_details.Patient_Name)
	}

	database.Postdb.QueryRow(context.Background(), "insert into prescriptions ($1,$2,$3,$4,$5)",
		Pres_details.Patient_Name, Pres_details.Medicine_name, Pres_details.Dosage_form, Pres_details.Quantity, Pres_details.Ph_num)

	
	c.JSON(200, gin.H{
		"message": "data added into database",
	})
}
