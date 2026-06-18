package notification_controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"lendogo-backend/internal/services"
)

type NotificationController struct {
	service services.NotificationService
}

func NewNotificationController(service services.NotificationService) *NotificationController {
	return &NotificationController{service: service}
}

func (c *NotificationController) GetUnread(ctx *fiber.Ctx) error {
	localUID, ok := ctx.Locals("user_id").(string)
	if !ok || localUID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userID, err := uuid.Parse(localUID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	notifs, err := c.service.GetUnreadNotifications(userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch notifications"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"notifications": notifs,
	})
}

func (c *NotificationController) MarkAllRead(ctx *fiber.Ctx) error {
	localUID, ok := ctx.Locals("user_id").(string)
	if !ok || localUID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userID, err := uuid.Parse(localUID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	err = c.service.MarkAsRead(userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to mark notifications as read"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
