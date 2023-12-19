package main

import (
	"io"

	"github.com/gin-gonic/gin"
)

func postEvent(c *gin.Context) {
	// print out request body for debugging:
	body, _ := io.ReadAll(c.Request.Body)
	c.Request.Body.Close()
	println(string(body))
}

func main() {
	router := gin.Default()
	router.POST("/events", postEvent)
	router.Run(":8080")
}
