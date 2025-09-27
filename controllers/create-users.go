package controllers

import (
	"PDS/database"
	"context"
	"github.com/gin-gonic/gin"
)

func CreateDocter(c *gin.Context) {

	var Docter struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`

	}
	if err := c.ShouldBindJSON(&Docter); err != nil{
		c.JSON(404 , gin.H{
			"error" : "Invalid input",
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
	res := database.Postdb.QueryRow(context.Background(), "SELECT CASE WHEN EXISTS ( SELECT 1 FROM users WHERE email = $1 AND role = 'admin') THEN 1 ELSE 0 END", email).Scan(&IsAdmin_db)

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
			
		})
		c.Abort()
		return
	}
	hashedPassword,err  := HashPassword(Docter.Password)
	if err!= nil{
		c.JSON(404 , gin.H{
			"error" : err , 
		})
	}
	database.Postdb.QueryRow(context.Background(), "insert into users (name , email , password , role) values ($1,$2,$3,$4)",Docter.Name , Docter.Email , hashedPassword , "docter" )

	

	c.JSON(200 , gin.H{
		"message" : "docter added successfully",
	})
}

func CreatePharmasist(c *gin.Context){

	var Pharmasist struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`

	}
	if err := c.ShouldBindJSON(&Pharmasist); err != nil{
		c.JSON(404 , gin.H{
			"error" : "Invalid input",
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
	res := database.Postdb.QueryRow(context.Background(), "SELECT CASE WHEN EXISTS ( SELECT 1 FROM users WHERE email = $1 AND role = 'admin') THEN 1 ELSE 0 END", email).Scan(&IsAdmin_db)

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
			
		})
		c.Abort()
		return
	}
	hashedPassword,err  := HashPassword(Pharmasist.Password)
	if err!= nil{
		c.JSON(404 , gin.H{
			"error" : err , 
		})
	}
	database.Postdb.QueryRow(context.Background(), "insert into users (name , email , password , role) values ($1,$2,$3,$4)",Pharmasist.Name , Pharmasist.Email , hashedPassword , "pharmasist" )

	c.JSON(200 , gin.H{
		"message" : "pharmasist added successfully",
	})
}

