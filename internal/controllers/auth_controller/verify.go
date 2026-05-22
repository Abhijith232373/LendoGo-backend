package auth

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	
	// Using the exact module name from your previous screenshot!
	"lendogo-backend/internal/database" 
)

// What we expect from React Step 2
type VerifyOTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func VerifyOTP(c *fiber.Ctx) error {
	var req VerifyOTPRequest

	// 1. Parse the JSON body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	ctx := context.Background()
	redisKey := "otp:" + req.Email

	// 2. Fetch the OTP from Redis
	storedOTP, err := database.RedisClient.Get(ctx, redisKey).Result()
	
	// If err is not nil, it means the OTP expired or the email is wrong!
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "OTP has expired or does not exist",
		})
	}

	// 3. Compare the user's OTP with the Redis OTP
	if storedOTP != req.OTP {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid OTP code",
		})
	}

	// 4. IT'S A MATCH! Delete the OTP from Redis so it can't be reused
	database.RedisClient.Del(ctx, redisKey)

	// Print success to your Docker terminal
	fmt.Printf("✅ SUCCESS: Email %s verified successfully!\n", req.Email)

	// 5. Tell React it was successful
	// Note: Later, we will generate a real Temporary JWT here!
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "OTP verified successfully",
		"tempToken": "fake-jwt-token-for-now", 
	})
}