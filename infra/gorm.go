package infra

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConnection *gorm.DB

func ConnectToDatabase() error {

	DSN := "user=%v password=%v dbname=%v port=%v sslmode=disable"

	config := postgres.Config{
		DSN:                  fmt.Sprintf(DSN, "root", "root", "loan_db", "5432"),
		PreferSimpleProtocol: true,
	}

	db, err := gorm.Open(postgres.New(config), &gorm.Config{})

	if err != nil {
		return err
	}

	DBConnection = db

	return nil
}
