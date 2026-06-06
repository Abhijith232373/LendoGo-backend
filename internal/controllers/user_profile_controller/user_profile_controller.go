package user_profile_controller

import (
	"fmt"
	"lendogo-backend/internal/services"
	"lendogo-backend/structures/dto"
	"lendogo-backend/structures/responses"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserProfileController struct {
	service services.UserProfileService
}

func NewUserProfileController(service services.UserProfileService) *UserProfileController {
	return &UserProfileController{service: service}
}

func (c *UserProfileController) GetProfile(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return responses.Error(ctx, 401, "Unauthorized")
	}

	data, err := c.service.GetMyProfile(userID)
	if err != nil {
		return responses.Error(ctx, 500, "Failed to load profile")
	}

	return responses.Success(ctx, 200, "Profile loaded", data)
}

func (c *UserProfileController) UpdateProfile(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return responses.Error(ctx, 401, "Unauthorized")
	}

	var req dto.UpdateProfileRequest
	// Use BodyParser which handles both JSON and Multipart Form Data
	if err := ctx.BodyParser(&req); err != nil {
		return responses.Error(ctx, 400, "Invalid request format")
	}

	// Handle Image Upload (Optional)
	imagePath := ""
	file, err := ctx.FormFile("profile_image")
	if err == nil {
		// Create a unique filename: user_id_timestamp.ext
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%s_%d%s", userID, time.Now().Unix(), ext)
		savePath := fmt.Sprintf("./uploads/profiles/%s", filename)
		
		if err := ctx.SaveFile(file, savePath); err == nil {
			// Save the URL path to the database
			imagePath = "/uploads/profiles/" + filename
		}
	}

	if err := c.service.UpdateMyProfile(userID, req, imagePath); err != nil {
		return responses.Error(ctx, 500, "Failed to update profile: "+err.Error())
	}

	return responses.Success(ctx, 200, "Profile updated successfully!", nil)
}