package routes

import (
	"github.com/gofiber/fiber/v2"
	
	// 👇 FIX 1: Import the exact admin_controller package you just created
	"lendogo-backend/internal/controllers/admin_controller" 
	"lendogo-backend/internal/middlewares"
)

// 👇 FIX 2: Point to admin_controller.AdminController
func SetupAdminRoutes(api fiber.Router, adminCtrl *admin_controller.AdminController) {
	// 1. Create a specific group for admin features
	adminGroup := api.Group("/admin")

	// 2. Apply BOTH middlewares to everything inside this group!
	adminGroup.Use(middlewares.Protected(), middlewares.AdminOnly())

	// 3. Add your admin-only routes
	adminGroup.Get("/all-users", adminCtrl.GetAllUsers)
	adminGroup.Get("/system-stats", adminCtrl.GetSystemStats)
}