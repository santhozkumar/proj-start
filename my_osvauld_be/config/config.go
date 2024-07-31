package config

import (
	"errors"

	"github.com/spf13/viper"
)

func GetJWTSecret() string {
	return viper.GetString("JWT_AUTH_SECRET")
}

func SetupConfig() {
	viper.AutomaticEnv()
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if secret := GetJWTSecret(); !validateJWTSecret(secret) {
		panic(errors.New("invalid jwt secret"))
	}
}

func validateJWTSecret(secret string) bool {
	if len(secret) != 32 {
		return false
	}
	return true
}
