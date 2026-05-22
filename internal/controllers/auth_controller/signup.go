package auth

import (
	"context"
	"fmt"

	// "lendogo-backen/internal/controllers/auth-controller"
	"lendogo-backend/internal/database"
	"lendogo-backend/internal/services/auth_service"

	// "lendogo-backend/internal/database"
	// "lendogo-backend/internal/services/auth-service"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
)

type OTPRequest struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
}

// Fiber requires handlers to take a *fiber.Ctx and return an error
func RequestOTP(c *fiber.Ctx) error {
	var req OTPRequest

	// 1. Parse the JSON body (Fiber makes this 1 line of code!)
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// 2. Generate a random 6-digit OTP
	rand.Seed(time.Now().UnixNano())
	otpCode := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 3. Save to Redis with a 5-minute expiration
	ctx := context.Background()
	redisKey := "otp:" + req.Email

	err := database.RedisClient.Set(ctx, redisKey, otpCode, 5*time.Minute).Err()
	if err != nil {
		fmt.Println("Redis Error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate OTP",
		})
	}

	// Call your new Email Service!
	// (Ensure you import your "lendogo-backen/internal/services/auth" package at the top)
	err = auth.SendOTPEmail(req.Email, otpCode)
	if err != nil {
		fmt.Println("SMTP Error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send email",
		})
	}

	// 4. Send success response back to React
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "OTP sent successfully",
	})
}