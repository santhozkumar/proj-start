package controller

import (
	"net/http"

	dto "osvauld/dtos"
	"osvauld/service"
	"osvauld/utils"

	"github.com/gin-gonic/gin"
)

func GetAuthor(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		// ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		SendResponse(ctx, http.StatusBadRequest, ResponseError, nil, "", 0, err)
		return
	}
	author, err := service.GetAuthor(ctx, userID)
	// ctx.AbortWithStatusJSON(http.StatusOK, author)
	SendResponse(ctx, http.StatusOK, ResponseSuccess, author, "fetched author", 0, nil)
}

func CreateAuthor(ctx *gin.Context) {
	var req dto.CreateAuthor
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, ResponseError, nil, "", 0, err)
		return
	}
	author, err := service.CreateAuthor(ctx, req)
	// ctx.AbortWithStatusJSON(http.StatusOK, author)
	SendResponse(ctx, http.StatusOK, ResponseSuccess, author, "fetched author", 0, nil)
}
