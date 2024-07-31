package controller

import (
	"log"
	"net/http"

	dto "osvauld/dtos"
	"osvauld/service"
	"osvauld/utils"

	"github.com/gin-gonic/gin"
)

func CreateFolder(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		// ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		SendResponse(ctx, http.StatusBadRequest, ResponseError, nil, "", 0, err)
		return
	}
	var req dto.CreateFolderRequest
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, ResponseError, nil, "", 0, err)
		return
	}
	folderRow, err := service.CreateFolder(ctx, userID, req)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, ResponseError, nil, "", 0, err)
		return
	}
	// ctx.AbortWithStatusJSON(http.StatusOK, author)
	log.Println("received create folder")
	SendResponse(ctx, http.StatusOK, ResponseSuccess, folderRow, "created folder", 0, nil)
}

func FetchFolderByUser(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		// ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		SendResponse(ctx, http.StatusBadRequest, ResponseError, nil, "", 0, err)
		return
	}

	folders, err := service.GetFoldersForUser(ctx, userID)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, ResponseError, nil, "", 0, err)
		return
	}
	// ctx.AbortWithStatusJSON(http.StatusOK, author)
	log.Println("received fetch folder for user")
	SendResponse(ctx, http.StatusOK, ResponseSuccess, folders, "fetched folders", 0, nil)
}
