package controllers

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Health(c *gin.Context) {
	// database.DBConnect()
	err := godotenv.Load()
	if os.Getenv("DATABASE_URL") == "" || err != nil {
		{
			c.JSON(500, gin.H{
				"error": "Database connection not configured",
			})
			return
		}

	}
	url := os.Getenv("DATABASE_URL")

	yes := strings.HasPrefix(url, "postgres://")
	if !yes {
		c.JSON(500, gin.H{
			"error": "Database connection not configured properly",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Health Is OK",
	})

}
