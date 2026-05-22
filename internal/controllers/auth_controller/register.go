package auth

import (
	"fmt"
	"lendogo-backend/internal/database"
	"lendogo-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	// Make sure this path perfectly matches your go.mod file!
	// "lendogo-backend/internal/database"
	// "lendogo-backend/internal/models"
)

type RegisterRequest struct {
    FullName string `json:"fullName"` 
    Email    string `json:"email"`
    Password string `json:"password"`
}

func FinalRegister(c *fiber.Ctx) error {
    var req RegisterRequest

    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request payload",
        })
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to hash password",
        })
    }

    user := models.User{
        Username:        req.FullName, 
        Email:           req.Email,
        Password:        string(hashedPassword),
        IsEmailVerified: true, 
    }

    if err := database.DB.Create(&user).Error; err != nil {
        fmt.Println("Database Insert Error:", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to save user to database",
        })
    }

    fmt.Println("✅ SUCCESS: User officially saved to PostgreSQL!")
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "User registered successfully",
    })
}