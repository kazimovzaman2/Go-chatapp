package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Config(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	return os.Getenv(key)
}

var JWTAccessSecret = Config("JWT_ACCESS_SECRET")
var JWTRefreshSecret = Config("JWT_REFRESH_SECRET")
