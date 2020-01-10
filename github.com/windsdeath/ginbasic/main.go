package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	r.POST("/form_post_with_querystring", formPost)

	return r
}

func formPost(c *gin.Context) {
	id := c.Query("id")
	strPage := c.DefaultQuery("page", "0")
	intPage, _ := strconv.Atoi(strPage)

	message := c.PostForm("message")
	nick := c.DefaultPostForm("nick", "anonymous")

	headerType := c.GetHeader("Content-Type")

	c.JSON(200, gin.H{
		"status":              "posted",
		"message":             message,
		"nick":                nick,
		"header-content-type": headerType,
		"id":                  id,
		"page":                intPage,
	})
}

func main() {
	v := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	v.Run(":8080")
}