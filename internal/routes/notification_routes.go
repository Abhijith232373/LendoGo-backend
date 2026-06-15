package routes

import (
	"github.com/gofiber/fiber/v2"

	"lendogo-backend/internal/controllers/notification_controller"
	"lendogo-backend/internal/middlewares"
)

func SetupNotificationRoutes(api fiber.Router, notifCtrl *notification_controller.NotificationController) {
	group := api.Group("/notifications")
	group.Use(middlewares.Protected())

	group.Get("/", notifCtrl.GetUnread)
	group.Post("/mark-read", notifCtrl.MarkAllRead)
}
