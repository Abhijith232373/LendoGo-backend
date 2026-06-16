package routes

import (
	"github.com/gofiber/fiber/v2"

	"lendogo-backend/internal/controllers/feedback_controller"
	"lendogo-backend/internal/middlewares"
	"lendogo-backend/internal/services"
)

func SetupFeedbackRoutes(api fiber.Router, feedbackCtrl *feedback_controller.FeedbackController, configService services.ConfigService) {
	feedbackGroup := api.Group("/feedback")

	// ==========================================
	// 🟢 PROTECTED USER ROUTES
	// ==========================================
	feedbackGroup.Post(
		"/", 
		middlewares.Protected(), // User must be logged in
		middlewares.RequireFeature(configService, "feedback"), // The Bouncer!
		feedbackCtrl.SubmitFeedback,
	)

	// ==========================================
	// 🔴 PROTECTED ADMIN ROUTES
	// ==========================================
	adminFeedbackGroup := feedbackGroup.Group("/admin")
	adminFeedbackGroup.Use(middlewares.Protected(), middlewares.AdminOnly())

	adminFeedbackGroup.Get("/", middlewares.RequirePermission("system.manage"), feedbackCtrl.GetAllFeedback)
	adminFeedbackGroup.Patch("/:id/status", middlewares.RequirePermission("system.manage"), feedbackCtrl.UpdateFeedbackStatus)
}