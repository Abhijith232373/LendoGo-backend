package feedback_controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"lendogo-backend/internal/services"
	"lendogo-backend/structures/dto"
	"lendogo-backend/structures/models"
)

type FeedbackController struct {
	feedbackService services.FeedbackService
}

func NewFeedbackController(fs services.FeedbackService) *FeedbackController {
	return &FeedbackController{feedbackService: fs}
}

func (c *FeedbackController) SubmitFeedback(ctx *fiber.Ctx) error {
	var req dto.SubmitFeedbackReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON payload"})
	}

	// 1. Basic validation: Rating must be within your UI bounds (e.g., 1 to 5)
	if req.Rating < 1 || req.Rating > 5 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Rating must be between 1 and 5"})
	}

	// 2. Safely extract and assert the User ID from the JWT context
	userIDRaw := ctx.Locals("user_id")
	if userIDRaw == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized access"})
	}

	userIDStr, ok := userIDRaw.(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user token format"})
	}

	// 3. Safely parse the UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Malformed user ID"})
	}

	// 4. Map to Model and send to Service
	feedback := models.Feedback{
		UserID:  userID,
		Rating:  req.Rating,
		Comment: req.Comment,
	}

	savedFeedback, err := c.feedbackService.SubmitFeedback(feedback)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Thank you! Your feedback helps us improve LendoGo.",
		"data":    savedFeedback,
	})
}

// GetAllFeedback handles Admin GET requests to fetch all platform feedback
func (c *FeedbackController) GetAllFeedback(ctx *fiber.Ctx) error {
	feedbacks, err := c.feedbackService.GetAllFeedback()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch feedback"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": feedbacks})
}

// UpdateFeedbackStatus handles Admin PATCH requests to change feedback status
func (c *FeedbackController) UpdateFeedbackStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	
	var req dto.UpdateFeedbackStatusReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON status format"})
	}

	err := c.feedbackService.UpdateStatus(id, req.Status)
	if err != nil {
		// Because we added business validation in the Service, we can return the exact error message here!
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Feedback status updated successfully"})
}