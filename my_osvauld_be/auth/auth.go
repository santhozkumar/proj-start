package auth

import (
	"osvauld/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claim struct {
	UserID   string `json:"token"`
	UserName string `json:"username"`
	jwt.RegisteredClaims
}


func GenerateJWT(userID, userName string) (string, error) {
    jwtSecret := config.GetJWTSecret()
	claim := Claim{
		userID,
        userName,
        jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        },
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
