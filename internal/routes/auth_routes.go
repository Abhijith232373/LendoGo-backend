package routes

import (
	authController "lendogo-backend/internal/controllers/auth_controller"

	"github.com/gofiber/fiber/v2"
)
func SetupAuthRoutes(router fiber.Router) {
    auth := router.Group("/auth")
    // Ensure these are POST and match the exact paths
    auth.Post("/request-otp", authController.RequestOTP)
    auth.Post("/verify-otp", authController.VerifyOTP)
    auth.Post("/register", authController.FinalRegister)
}