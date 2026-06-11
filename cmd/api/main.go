package main

import "github.com/gin-gonic/gin"

func main() {
	var router *gin.Engine = gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			//map [string] of interface
			//map [string] any {}
			"message": "Todo API is running!!!!",
			"status":  "success",
		})
	})
	router.Run(":8080")

}
