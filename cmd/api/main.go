package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv" 

	"lendogo-backend/internal/database" // <-- 1. This imports your database recipe!
	"lendogo-backend/internal/routes"
)

func main() {
	// 1. Force Go to read the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️ Warning: Could not load .env file")
	} else {
		fmt.Println("✅ Loaded .env file successfully!")
	}

	// 2. Boot up Redis (for OTPs)
	database.InitRedis()

	// 3. 👉 BOOT UP POSTGRES (This forces AutoMigrate to create the users table!)
	err = database.Connect()
	if err != nil {
		log.Fatalf("❌ Database Boot Failed: %v\n", err)
	}

	// 4. Initialize Fiber API
	app := fiber.New()

	// CORS Setup
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Setup Routes
	routes.SetupAuthRoutes(app)

	fmt.Println("🚀 Fiber Server running on port 8080...")
	log.Fatal(app.Listen("0.0.0.0:8080"))
}