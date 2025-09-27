package controllers

import (
	db "PDS/database"
	"PDS/midlewares"
	"PDS/models"
	"context"

	"net/http"

	"github.com/gin-gonic/gin"
)

func Add_Medicine(c *gin.Context) {
	var medicine models.Add_Medicine

	if err := c.ShouldBindJSON(&medicine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	emailVal, exists_1 := c.Get("email")
    if !exists_1 {
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

	isadmin := db.Postdb.QueryRow(context.Background(), "SELECT CASE WHEN EXISTS ( SELECT 1 FROM users WHERE email = $1 AND role = 'admin') THEN 1 ELSE 0 END", email).Scan(&IsAdmin_db)

	if isadmin != nil {
		c.JSON(404, gin.H{
			"error": isadmin,
		})
	}

	if IsAdmin_db == 0 {
		c.JSON(400, gin.H{
			"message": "you are not allowed to do this functions",
			"IsAdmin": isadmin,
		})
		c.Abort()
		return
	}

	var exists bool

	checkQuery := `SELECT EXISTS(SELECT 1 FROM medicine WHERE medicine_name=$1)`
	err := db.Postdb.QueryRow(context.Background(), checkQuery, medicine.Medicine_Name).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if exists {

		updateQuery := `UPDATE medicine SET stock_quantity = stock_quantity + $1 
						WHERE medicine_name=$2 
						RETURNING stock_quantity`
		err := db.Postdb.QueryRow(context.Background(), updateQuery, medicine.Stock_Quantity, medicine.Medicine_Name).
			Scan(&medicine.Stock_Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":        "stock updated successfully",
			"medicine_name":  medicine.Medicine_Name,
			"stock_quantity": medicine.Stock_Quantity,
		})
		return
	}

	insertQuery := `INSERT INTO medicine (medicine_name, dosage_form, stock_quantity)
					VALUES ($1, $2, $3)
					RETURNING medicine_name, dosage_form, stock_quantity`
	err = db.Postdb.QueryRow(context.Background(), insertQuery,
		medicine.Medicine_Name, medicine.Dosage_form, medicine.Stock_Quantity).
		Scan(&medicine.Medicine_Name, &medicine.Dosage_form, &medicine.Stock_Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "medicine added successfully",
		"medicine_name":  medicine.Medicine_Name,
		"dosage_form":    medicine.Dosage_form,
		"stock_quantity": medicine.Stock_Quantity,
	})
}

func CreateAdmin(c *gin.Context) {
	var Admin struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&Admin); err != nil {
		c.JSON(404, gin.H{
			"errror": "invalid input",
		})
		c.Abort()
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

	isadmin := db.Postdb.QueryRow(context.Background(), "SELECT CASE WHEN EXISTS ( SELECT 1 FROM users WHERE email = $1 AND role = 'admin') THEN 1 ELSE 0 END", email).Scan(&IsAdmin_db)

	if isadmin != nil {
		c.JSON(404, gin.H{
			"error": isadmin,
		})
	}

	if IsAdmin_db == 0 {
		c.JSON(400, gin.H{
			"message": "you are not allowed to do this functions",
			"IsAdmin": isadmin,
		})
		c.Abort()
		return
	}
	hashpassword, err := HashPassword(Admin.Password)
	if err != nil {
		c.JSON(404, gin.H{
			"error": "error in hashing password ",
		})
	}
	_, res := db.Postdb.Exec(context.Background(), "insert into users (name , email , password , role)  values ($1,$2,$3,$4)", Admin.Name, Admin.Email, hashpassword, "admin")
	if res != nil {
		c.JSON(404, gin.H{
			"error": res,
			"pass":  hashpassword,
		})
	}
	token, err := midlewares.GenerateToken(Admin.Email)
	if err != nil {
		c.JSON(404, gin.H{
			"error": err,
		})
	}

	c.JSON(200, gin.H{
		"message": "successfully",
		"pass":    hashpassword,
		"admin":   email,
		"token":   token,
	})
}
