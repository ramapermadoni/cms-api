package main

import (
	"cms-api/config"
	"cms-api/internal/categories"
	"cms-api/internal/database/connection"
	"cms-api/internal/database/migration"
	"cms-api/internal/media"
	"cms-api/internal/posts"
	"cms-api/internal/users"
	"cms-api/pkg/utility/logger"
	"cms-api/seeders"
	"flag"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// seed := flag.Bool("seed", false, "Run database seeder")
	flag.Parse()
	config.Initiator()
	logger.Initiator()
	connection.InitDB()
	migration.AutoMigrate()
	// Jalankan seeding jika flag --seed digunakan
	seed := viper.GetBool("seed")
	if seed {
		fmt.Println("Seeding database...")
		seeders.SeedAll(connection.DB)
		return // Hentikan aplikasi setelah seeding
	} else {
		fmt.Println("Not seeding database...")
	}
	InitiateRouter()
}

func InitiateRouter() {
	fmt.Println("app mode: ", viper.GetString("app.mode"))
	mode := viper.GetString("app.mode")
	if mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://192.168.180.27:3000", // Origin untuk Next.js development
			"http://localhost:3000",      // Origin lain untuk local development
			"https://yourproduction.com", // Origin untuk production
		}, // or "*" for all origins
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// // Configure CORS
	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:3000"}, // or "*" for all origins
	// 	AllowMethods:     []string{"GET", "POST", "OPTIONS"},
	// 	AllowHeaders:     []string{"Content-Type", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))
	// Route default untuk root ("/")
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":       "Welcome to CMS API!",
			"repository":    "https://github.com/ramapermadoni/cms-api",
			"readme":        "https://github.com/ramapermadoni/cms-api/#readme",
			"documentation": "https://documenter.getpostman.com/view/15292179/2sAXxY2T1J",
		})
	})
	users.Initiator(router)
	users.AuthInitiator(router)
	categories.Initiator(router)
	posts.Initiator(router)
	media.Initiator(router)

	// Serve folder "uploads" agar bisa diakses publik
	router.Static("/uploads", "./uploads")
	router.Run(":8080")
}
