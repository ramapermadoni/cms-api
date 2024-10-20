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

	users.Initiator(router)
	users.AuthInitiator(router)
	categories.Initiator(router)
	posts.Initiator(router)
	media.Initiator(router)

	// Serve folder "uploads" agar bisa diakses publik
	router.Static("/uploads", "./uploads")
	router.Run(":8080")
}
