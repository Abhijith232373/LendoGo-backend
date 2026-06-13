package career_controller

import (

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"lendogo-backend/internal/services"
	"lendogo-backend/structures/models"
	"lendogo-backend/utils"
)

type CareerController struct {
	careerService services.CareerService
}

func NewCareerController(cs services.CareerService) *CareerController {
	return &CareerController{careerService: cs}
}

// ==========================================
// JOB OPENING METHODS (ADMIN & PUBLIC)
// ==========================================

// CreateOpening handles POST requests from the Admin HR Panel
func (c *CareerController) CreateOpening(ctx *fiber.Ctx) error {
	var req models.CareerOpening
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format. Check array fields."})
	}

	opening, err := c.careerService.CreateOpening(req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Job opening published successfully",
		"data":    opening,
	})
}

// GetOpenings handles GET requests for the public careers page
func (c *CareerController) GetOpenings(ctx *fiber.Ctx) error {
	// React frontend can send ?status=Open to only see active jobs
	status := ctx.Query("status")

	openings, err := c.careerService.GetAllOpenings(status)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch current openings"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": openings})
}

// GetOpeningByID handles GET requests for the job detail page (Read More)
func (c *CareerController) GetOpeningByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	opening, err := c.careerService.GetOpeningByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Job opening not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": opening})
}

// UpdateOpening handles PUT requests to edit an existing job opening
func (c *CareerController) UpdateOpening(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req models.CareerOpening
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format."})
	}
	opening, err := c.careerService.UpdateOpening(id, req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Job opening updated successfully", "data": opening})
}

// UpdateOpeningStatus handles PATCH requests to toggle job status
func (c *CareerController) UpdateOpeningStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format."})
	}
	err := c.careerService.UpdateOpeningStatus(id, req.Status)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Status updated"})
}

// ==========================================
// JOB APPLICATION METHOD (FILE UPLOAD)
// ==========================================

// SubmitApplication handles multipart/form-data for job applications + resumes
func (c *CareerController) SubmitApplication(ctx *fiber.Ctx) error {
	jobIDStr := ctx.Params("id")
	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid job opening ID"})
	}

	// 1. Grab the Resume File from the request
	file, err := ctx.FormFile("resume")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Resume file is required"})
	}

	// 2. Validate File Size (Max 5MB as requested in your UI)
	if file.Size > 5*1024*1024 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Resume must be less than 5MB"})
	}

	// 3. Upload to S3
	s3URL, uploadErr := utils.UploadFileToS3(file)
	if uploadErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload resume to S3"})
	}

	// 4. Build the Model from the Form Data
	application := models.JobApplication{
		CareerOpeningID: jobID,
		FirstName:       ctx.FormValue("first_name"),
		LastName:        ctx.FormValue("last_name"),
		Email:           ctx.FormValue("email"),
		Phone:           ctx.FormValue("phone"),
		Address:         ctx.FormValue("address"),
		City:            ctx.FormValue("city"),
		State:           ctx.FormValue("state"),
		PostalCode:      ctx.FormValue("postal_code"),
		ResumePath:      s3URL, // Store the S3 URL
	}

	// Basic validation
	if application.FirstName == "" || application.Email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "First Name and Email are required fields"})
	}

	// 5. Send to Service
	savedApp, err := c.careerService.SubmitApplication(application)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Application submitted successfully! Our HR team will reach out soon.",
		"data":    savedApp,
	})
}

// GetAllApplications handles GET requests to fetch all candidate applications
func (c *CareerController) GetAllApplications(ctx *fiber.Ctx) error {
	applications, err := c.careerService.GetAllApplications()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch job applications"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": applications})
}

// UpdateApplicationStatus handles PATCH requests to update a candidate's application status
func (c *CareerController) UpdateApplicationStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format."})
	}
	
	err := c.careerService.UpdateApplicationStatus(id, req.Status)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Application status updated successfully"})
}