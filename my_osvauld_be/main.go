package main

import (
	"fmt"
	"log"
	"net/http"
	"osvauld/auth"
	"osvauld/config"
	"osvauld/controller"
	"osvauld/infra/database"
	"osvauld/router/middleware"

	"github.com/gin-contrib/requestid"

	"github.com/gin-gonic/gin"
)

func test(c *gin.Context) {
	c.Writer.Header().Set("Location", "https://www.google.com")
	c.AbortWithStatusJSON(http.StatusFound, gin.H{"success": true})
}

func create(c *gin.Context) {
	userID := "7648fea8-72d7-4f58-8b05-334d52dce19e"
	userName := "admin"
	token, err := auth.GenerateJWT(userID, userName)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"success": false, "error": fmt.Sprintf("Generate JWT failed: %s\n", err)})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"success": true, "token": token})
}

func main() {
	config.SetupConfig()
	masterDSN, _ := config.DbConfiguration()
	database.DBConnection(masterDSN)

	server := gin.New()
	server.Use(requestid.New(), middleware.CorsMiddleware(), gin.Logger())

	privateApi := server.Group("/")
	privateApi.Use(middleware.AuthMiddleWare())
	privateApi.GET("/admin", test)

	privateApi.GET("/author", controller.GetAuthor)
	privateApi.POST("/author", controller.CreateAuthor)
	privateApi.GET("/folders", controller.FetchFolderByUser)
	privateApi.POST("/folders", controller.CreateFolder)

	publicApi := server.Group("/")
	publicApi.GET("/create", create)
	publicApi.GET("/test", test)

	server.Use(gin.Recovery())
	log.Fatal(server.Run(":8080"))
}
