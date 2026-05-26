package admin_controller

import (
	"github.com/gofiber/fiber/v2"
)

// AdminController structure
type AdminController struct {
	// Later, we will inject your Database Repositories or Services here!
}

// Constructor for the Controller
func NewAdminController() *AdminController {
	return &AdminController{}
}

// ==========================================
// ADMIN ROUTES
// ==========================================

// GetAllUsers - Only admins can see this!
func (c *AdminController) GetAllUsers(ctx *fiber.Ctx) error {
	// For now, just return a success message. Later, fetch users from DB.
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success! Here is the list of all users. (Admin Eyes Only)",
	})
}

// GetSystemStats - Only admins can see this!
func (c *AdminController) GetSystemStats(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "System is running perfectly.",
		"active_loans": 42,
	})
}