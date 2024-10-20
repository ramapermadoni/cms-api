package seeders

import (
	"cms-api/internal/categories"
	"log"

	"gorm.io/gorm"
)

func SeedCategories(db *gorm.DB) {

	categories := []categories.Category{
		{Name: "News", Description: "Latest updates and important events to keep users informed."},
		{Name: "Articles", Description: "Informative content that provides in-depth discussion on specific topics."},
		{Name: "Blog", Description: "Casual posts featuring opinions, stories, or personal and business experiences."},
		{Name: "Static Pages", Description: "Permanent and essential information, such as \"About Us\" or \"Contact.\""},
		{Name: "Gallery", Description: "A collection of images or videos showcasing visual documentation."},
		{Name: "Products", Description: "A catalog of goods or services offered, including descriptions and prices."},
		{Name: "FAQ", Description: "Answers to common questions to help users understand the service."},
		{Name: "Events", Description: "Information about upcoming or past events with relevant details."},
		{Name: "Portfolio", Description: "A showcase of projects or work highlighting skills and achievements."},
		{Name: "Announcements", Description: "Brief notifications about urgent or important information."},
		{Name: "Documents", Description: "Files and documents available for users to download or access."},
		{Name: "Testimonials", Description: "Customer reviews that build credibility for products or services."},
	}

	for _, category := range categories {
		if err := db.Create(&category).Error; err != nil {
			log.Fatalf("Gagal membuat category %s: %v", category.Name, err)
		}
	}

	log.Println("Berhasil menambahkan data kategori.")
}
