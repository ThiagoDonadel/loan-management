package infra

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConnection *gorm.DB

func ConnectToDatabase() error {

	DSN := fmt.Sprintf(
		"user=%v password=%v dbname=%v port=%v sslmode=disable",
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.dbname"),
		viper.GetString("database.port"),
	)

	config := postgres.Config{
		DSN:                  DSN,
		PreferSimpleProtocol: true,
	}

	db, err := gorm.Open(postgres.New(config), &gorm.Config{})

	if err != nil {
		return err
	}

	DBConnection = db

	return nil
}
