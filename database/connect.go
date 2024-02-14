package database

import (
	"fmt"
	"strconv"

	"github.com/kazimovzaman2/Go-chatapp/config"
	"github.com/kazimovzaman2/Go-chatapp/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(config *config.Config) {
	var err error
	p := config.DBPort
	port, err := strconv.ParseUint(p, 10, 32)

	if err != nil {
		panic("Failed to parse database port.")
	}

	dns := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		port,
		config.DBUserName,
		config.DBUserPassword,
		config.DBName,
	)
	DB, err = gorm.Open(postgres.Open(dns), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	DB.AutoMigrate(&model.User{})
	fmt.Println("✅ Database connected.")
}
