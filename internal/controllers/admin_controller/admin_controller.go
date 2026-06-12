package admin_controller

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"lendogo-backend/database"
	"lendogo-backend/internal/services" 
	"lendogo-backend/structures/models"
	
	// 👇 THE MAGIC IMPORT
	"lendogo-backend/internal/websockets" 
)

// AdminController structure handles administrative HTTP requests.
type AdminController struct {
	adminService services.AdminService 
}

// NewAdminController initializes a new AdminController.
func NewAdminController(as services.AdminService) *AdminController {
	return &AdminController{adminService: as}
}

// ==========================================
// 0. AUTHENTICATION (Staff Login)
// ==========================================

type AdminLoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *AdminController) AdminLogin(ctx *fiber.Ctx) error {
	var req AdminLoginReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	staff, err := c.adminService.AdminLogin(req.Email, req.Password)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	claims := jwt.MapClaims{
		"user_id": staff.ID.String(),
		"role":    "admin", 
		"exp":     time.Now().Add(time.Hour * 24).Unix(), 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "my_super_secret_lendo_go_key_998877" 
	}
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not log in"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   t,
		"staff": fiber.Map{
			"id":          staff.ID,
			"full_name":   staff.FullName,
			"email":       staff.Email,
			"role":        staff.Role,
			"permissions": staff.Permissions, 
		},
	})
}

// ==========================================
// 1. STAFF MANAGEMENT
// ==========================================

func (c *AdminController) CreateStaff(ctx *fiber.Ctx) error {
	var req services.CreateStaffDTO
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid form data"})
	}

	if err := c.adminService.CreateStaff(req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to provision staff account"})
	}

	// 🔴 WEBSOCKET BROADCAST
	websockets.BroadcastMessage("STAFF_PROVISIONED", fiber.Map{
		"message": "A new internal staff account was created.",
		"email":   req.Email,
	})

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Staff account provisioned successfully!"})
}

func (c *AdminController) GetAllStaff(ctx *fiber.Ctx) error {
	staff, err := c.adminService.GetAllStaff()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch staff directory"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Staff directory fetched", "data": staff})
}

// ==========================================
// 2. USER MANAGEMENT
// ==========================================

func (c *AdminController) GetAllUsers(ctx *fiber.Ctx) error {
	var users []models.User
	result := database.DB.Omit("password").Preload("Profile").Order("created_at DESC").Find(&users)
	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Users fetched", "data": users})
}

func (c *AdminController) CreateUser(ctx *fiber.Ctx) error {
	// ... (Your existing struct and body parser code) ...
	var req struct {
		FullName     string `json:"full_name"`
		Email        string `json:"email"`
		Role         string `json:"role"`
		MobileNumber string `json:"mobile_number"`
		DOB          string `json:"dob"`
		PanCard      string `json:"pan_card_number"`
		CreditRating string `json:"credit_rating"`
		CreditScore  int    `json:"credit_score"`
		Address      string `json:"address"`
		City         string `json:"city"`
		State        string `json:"state"`
		Pincode      string `json:"pincode"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload format"})
	}

	b := make([]byte, 4)
	_, _ = rand.Read(b)
	plainPassword := "Lendo" + hex.EncodeToString(b)[:4] + "@"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to secure key"})
	}

	userID := uuid.New()

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		userRecord := models.User{
			ID:              userID,
			FullName:        req.FullName,
			Email:           req.Email,
			Password:        string(hashedPassword),
			Role:            req.Role,
			IsEmailVerified: true,
			Status:          "Active",
		}
		if err := tx.Create(&userRecord).Error; err != nil { return err }

		profileRecord := models.UserProfile{
			UserID:        userID,
			PhoneNumber:   req.MobileNumber,
			DateOfBirth:   req.DOB,
			Address:       req.Address,
			City:          req.City,
			State:         req.State,
			Pincode:       req.Pincode,
			TrustScore:    req.CreditScore,
			PanCardNumber: req.PanCard,
			CreditRating:  req.CreditRating,
		}
		if err := tx.Create(&profileRecord).Error; err != nil { return err }
		return nil
	})

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB crash"})
	}

	// 🔴 WEBSOCKET BROADCAST
	websockets.BroadcastMessage("USER_CREATED", fiber.Map{
		"user_id":   userID,
		"full_name": req.FullName,
		"email":     req.Email,
	})

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully!",
		"default_password": plainPassword,
	})
}

func (c *AdminController) UpdateUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var payload map[string]interface{}
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	updates := map[string]interface{}{
		"full_name": payload["full_name"],
		"email":     payload["email"],
		"role":      payload["role"],
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	// ... (Your existing profile updates code) ...

	// 🔴 WEBSOCKET BROADCAST
	websockets.BroadcastMessage("USER_UPDATED", fiber.Map{
		"user_id": id,
		"message": "A user profile was updated.",
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User updated"})
}

func (c *AdminController) DeleteUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	result := database.DB.Where("id = ?", id).Delete(&models.User{})
	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database crash"})
	}
	if result.RowsAffected == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}

	// 🔴 WEBSOCKET BROADCAST
	websockets.BroadcastMessage("USER_DELETED", fiber.Map{
		"user_id": id,
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted"})
}

func (c *AdminController) UpdateUserStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var payload struct { Status string `json:"status"` }
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	database.DB.Exec("UPDATE users SET status = ? WHERE id = ?", payload.Status, id)

	// 🔴 WEBSOCKET BROADCAST
	websockets.BroadcastMessage("USER_STATUS_UPDATED", fiber.Map{
		"user_id": id,
		"status":  payload.Status,
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Status updated"})
}

// ==========================================
// 3. LOANS & SYSTEM DASHBOARD
// ==========================================

func (c *AdminController) GetSystemStats(ctx *fiber.Ctx) error {
	// GET route - no broadcast needed
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "System is running", "active_loans": 42})
}

func (c *AdminController) GetAllApplications(ctx *fiber.Ctx) error {
	// GET route - no broadcast needed
	var applications []models.LoanApplication
	database.DB.Preload("KYC").Preload("FinancialDocs").Order("created_at DESC").Find(&applications)
	// ... (Your presigned URL generation logic) ...
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": applications})
}

func (c *AdminController) UpdateApplicationStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var payload struct { Status string `json:"status"` }
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "malformed JSON payload"})
	}

	validStates := map[string]bool{
		"APPROVED": true, "REJECTED": true, "ADDITIONAL_DOCS_REQUIRED": true, "DISBURSED": true,
	}

	if !validStates[payload.Status] {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid state"})
	}

	if payload.Status == "DISBURSED" {
		err := database.DB.Transaction(func(tx *gorm.DB) error {
			// ... (Your existing disbursement transaction logic) ...
			var loan models.LoanApplication
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&loan).Error; err != nil { return err }
			if loan.ApplicationStatus == "DISBURSED" { return nil }
			
			var sysWallet models.SystemWallet
			tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("wallet_name = ?", "capital_disbursement").First(&sysWallet)
			tx.Model(&sysWallet).UpdateColumn("balance", gorm.Expr("balance - ?", loan.PrincipalAmount))

			var userWallet models.UserWallet
			tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", loan.UserID).FirstOrCreate(&userWallet, models.UserWallet{UserID: loan.UserID, Balance: 0})
			tx.Model(&userWallet).UpdateColumn("balance", gorm.Expr("balance + ?", loan.PrincipalAmount))
			
			return tx.Model(&loan).UpdateColumn("application_status", "DISBURSED").Error
		})

		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// 🔴 WEBSOCKET BROADCAST FOR DISBURSEMENT
		websockets.BroadcastMessage("LOAN_DISBURSED", fiber.Map{
			"loan_id": id,
			"message": "Capital has been moved to user wallet.",
		})

		return ctx.SendStatus(fiber.StatusOK)
	}

	result := database.DB.Model(&models.LoanApplication{}).Where("id = ?", id).Update("application_status", payload.Status)
	if result.Error != nil || result.RowsAffected == 0 {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed update"})
	}

	// 🔴 WEBSOCKET BROADCAST FOR REGULAR STATUS UPDATE
	websockets.BroadcastMessage("LOAN_STATUS_UPDATED", fiber.Map{
		"loan_id": id,
		"status":  payload.Status,
	})

	return ctx.SendStatus(fiber.StatusOK)
}

func (c *AdminController) GetAllConsultations(ctx *fiber.Ctx) error {
	// GET route - no broadcast needed
	var consultations []models.Consultation
	database.DB.Order("created_at DESC").Find(&consultations)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": consultations})
}