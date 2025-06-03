package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/zirris512/Authentication/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true

	router.Use(cors.New(config))

	router.POST("/authentication/register", func(c *gin.Context) {
		var postBody PostBody
		if err := c.ShouldBindJSON(&postBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newUser, err := Register(postBody.Username, postBody.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = database.AddToken(newUser.Username, newUser.Token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, newUser)
	})

	router.POST("/authentication/login", func(c *gin.Context) {
		var postBody PostBody
		if err := c.ShouldBindJSON(&postBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		verifiedUser, err := Login(postBody.Username, postBody.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		err = database.AddToken(verifiedUser.Username, verifiedUser.Token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, verifiedUser)
	})

	router.POST("/authentication/verify", func(c *gin.Context) {
		var postBodyToken struct {
			Token string `json:"token"`
		}
		if err := c.ShouldBindJSON(&postBodyToken); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := VerifyToken(postBodyToken.Token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}
		user, err := database.GetUserByToken(postBodyToken.Token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"username":  user.Username,
			"watchlist": user.Watchlist,
		})
	})

	router.PUT("/authentication/watchlist", func(c *gin.Context) {
		var userWatchlist UserWatchlist
		if err := c.ShouldBindJSON(&userWatchlist); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := database.UpdateWatchlist(userWatchlist.Username, userWatchlist.Watchlist)
		if err != nil {
			c.Status(http.StatusBadRequest)
		}
		c.Status(http.StatusCreated)
	})

	router.Run(":3301")
}
