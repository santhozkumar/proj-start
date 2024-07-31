package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserIDFromContext(ctx *gin.Context) (uuid.UUID, error) {
    userIDInterface, _ := ctx.Get("userID")
    fmt.Println(ctx.Keys)
    fmt.Println(userIDInterface)
    userID, ok := userIDInterface.(uuid.UUID)
    fmt.Println(userID)
	if !ok{
		return uuid.UUID{}, fmt.Errorf("Userid %s not satifsfying the UUID Interface", userID)
	}
    return userID, nil
}
