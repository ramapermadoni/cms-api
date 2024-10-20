package migration

import (
	"cms-api/internal/categories"
	"cms-api/internal/database/connection"
	"cms-api/internal/media"
	"cms-api/internal/posts"
	"cms-api/internal/users"
	"log"
)

// autoMigrate migrates the defined models
func AutoMigrate() {
	err := connection.DB.AutoMigrate(
		&users.User{},
		&categories.Category{},
		&posts.Post{},
		&media.Media{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrated successfully!")
}
