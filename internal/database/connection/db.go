package connection

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

// InitDB initializes the database connection and migrates the schema
func InitDB() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("migration.db.postgres.db_host"),
		viper.GetInt("migration.db.postgres.db_port"),
		viper.GetString("migration.db.postgres.db_user"),
		viper.GetString("migration.db.postgres.db_password"),
		viper.GetString("migration.db.postgres.db_name"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Optionally, you can perform a connection test here if needed
	fmt.Println("Successfully connected to the database")
}
