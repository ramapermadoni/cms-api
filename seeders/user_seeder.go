package seeders

import (
	"cms-api/internal/users"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	password, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)

	user := users.User{
		Username: "admin",
		Fullname: "Administrator CMS",
		Email:    "admin@cmsapi.com",
		Password: string(password),
		Role:     "admin",
	}

	if err := db.Create(&user).Error; err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}
	log.Println("Berhasil menambahkan data user.")
}
