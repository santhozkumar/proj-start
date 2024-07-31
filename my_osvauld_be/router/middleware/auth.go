package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"osvauld/auth"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("called authorization")
		header := c.GetHeader("Authorization")
		if header == "" {
			log.Println("no header return")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
			return
		}
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Token should start with Bearer"})
			return
		}
		tokenString := strings.TrimPrefix(header, "Bearer ")

		keyfunc := func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return []byte(""), errors.New("Invalid singing method")

			}
			log.Println("Headers from jwt", token.Header)
			key := "6M6H5u8DJnWxg33bgcpGaLs6k4pAE7x9"
			return []byte(key), nil
		}
		// Keyfunc func(*Token) (interface{}, error)

		claim := &auth.Claim{}
		token, err := jwt.ParseWithClaims(tokenString, claim, keyfunc, jwt.WithValidMethods([]string{"HS256"}))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error parsing the token:%s\n", err)})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			return
		}

		userID, err := uuid.Parse(claim.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID token"})
			return
		}

		// c.Header("userID", userID.String())
        c.Set("userID", userID)
        log.Println(userID, "user id set from the server")
		c.Next()
	}
}
