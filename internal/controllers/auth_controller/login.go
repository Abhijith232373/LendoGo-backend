package auth_controller

import (
	"lendogo-backend/structures/dto"

	"github.com/gofiber/fiber/v2"
)

// The Login Handler
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req dto.LoginReq

	// Step 1: Read the JSON (email & password) from the React frontend
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Step 2: Pass the request to the Service layer we just built
	res, err := c.authService.Login(req)
	if err != nil {
		// If passwords don't match or email isn't found, return 401 Unauthorized
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Step 3: Success! Send the Token and User Data back to React
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"data":    res, // This contains the token and user details!
	})
}