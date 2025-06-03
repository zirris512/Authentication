package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true

	router.Use(cors.New(config))

	router.POST("/authentication/register", func(c *gin.Context) {
		var postBody PostBody
		if err := c.ShouldBindJSON(&postBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		newUser := Register(postBody.Username, postBody.Password)
		c.JSON(http.StatusOK, newUser)
	})

	router.POST("/authentication/login", func(c *gin.Context) {
		var postBody PostBody
		if err := c.ShouldBindJSON(&postBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		verifiedUser := Login(postBody.Username, postBody.Password)
		c.JSON(http.StatusOK, verifiedUser)
	})

	router.POST("/authentication/verify", func(c *gin.Context) {
		var postBodyToken struct {
			Token string `json:"token"`
		}
		if err := c.ShouldBindJSON(&postBodyToken); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		err := VerifyToken(postBodyToken.Token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"authorized": true,
		})
	})

	router.Run(":3301")
}
