package main

import (
    "log"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/joho/godotenv" 
    
    "lendogo-backend/database"
    "lendogo-backend/internal/controllers/auth_controller"
    "lendogo-backend/internal/repositories"
    "lendogo-backend/internal/routes"
    "lendogo-backend/internal/services"
)

func main() {
    // 1. Load the .env file before doing ANYTHING else
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: No .env file found or failed to load")
    }

    app := fiber.New()

    // 2. Enable CORS so React (localhost:5173) can talk to Go (localhost:8080)
    app.Use(cors.New(cors.Config{
        AllowOrigins: "http://localhost:5173",
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    }))

    // ==========================================
    // 3. DATABASE CONNECTIONS
    // ==========================================
    
    // Connect to PostgreSQL
    err := database.Connect()
    if err != nil {
        log.Fatal("Failed to connect to PostgreSQL:", err)
    }

    // ADDED THIS: Connect to Redis!
    err = database.ConnectRedis()
    if err != nil {
        log.Fatal("Failed to connect to Redis:", err)
    }

    // ==========================================
    // 4. DEPENDENCY INJECTION WIRING
    // ==========================================
    
    // Step A: Give the Global Database to the Repository
    userRepo := repositories.NewUserRepository(database.DB)
    
    // Step B: Give the Repository to the Service
    authService := services.NewAuthService(userRepo)
    
    // Step C: Give the Service to the Controller
    authController := auth_controller.NewAuthController(authService)

    // ==========================================
    // 5. Setup Routes
    // ==========================================
    api := app.Group("/api")
    routes.SetupAuthRoutes(api, authController)

    // 6. Start the Server
    log.Println("🚀 Fiber Server running on port 8080...")
    log.Fatal(app.Listen(":8080"))
}	