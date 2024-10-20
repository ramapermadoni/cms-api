package seeders

import "gorm.io/gorm"

func SeedAll(db *gorm.DB) {
	SeedUsers(db)
	SeedCategories(db)
}
