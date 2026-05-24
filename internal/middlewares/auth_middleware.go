package middlewares

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected route middleware to verify JWT tokens
func Protected() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Step 1: Get the Authorization header from the request
		authHeader := ctx.Get("Authorization")

		// Step 2: Check if the header is missing or doesn't start with "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Missing or invalid token format",
			})
		}

		// Step 3: Extract the token string (remove "Bearer " from the front)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Step 4: Parse and validate the token
		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is what we expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		// Step 5: Handle invalid tokens or errors
		if err != nil || !token.Valid {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid or expired token",
			})
		}

		// Step 6: Extract the user_id from the token claims and store it in the context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Save the user ID into Fiber's local storage so the next Controller can use it
			ctx.Locals("user_id", claims["user_id"])
			
			// SUCCESS! Pass the request to the actual Controller
			return ctx.Next()
		}

		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: Failed to process token claims",
		})
	}
}