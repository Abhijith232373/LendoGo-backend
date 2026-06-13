package repositories

import (
	"lendogo-backend/structures/models"

	"gorm.io/gorm"
)

type CareerRepository interface {
	// Job Opening Management
	CreateOpening(opening *models.CareerOpening) error
	GetAllOpenings(status string) ([]models.CareerOpening, error)
	GetOpeningByID(id string) (*models.CareerOpening, error)
	UpdateOpening(id string, opening *models.CareerOpening) error
	UpdateOpeningStatus(id string, status string) error
	
	// 👇 The missing method to save candidate applications!
	SubmitApplication(application *models.JobApplication) error
	GetAllApplications() ([]models.JobApplication, error)
	UpdateApplicationStatus(id string, status string) error
}

type careerRepository struct {
	db *gorm.DB
}

// NewCareerRepository creates a new instance of the repository
func NewCareerRepository(db *gorm.DB) CareerRepository {
	return &careerRepository{db: db}
}

// ==========================================
// JOB OPENING METHODS
// ==========================================

func (r *careerRepository) CreateOpening(opening *models.CareerOpening) error {
	return r.db.Create(opening).Error
}

func (r *careerRepository) GetAllOpenings(status string) ([]models.CareerOpening, error) {
	var openings []models.CareerOpening
	query := r.db.Model(&models.CareerOpening{})
	
	// If a status is provided (like "Open"), filter by it
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Order("created_at DESC").Find(&openings).Error
	return openings, err
}

func (r *careerRepository) GetOpeningByID(id string) (*models.CareerOpening, error) {
	var opening models.CareerOpening
	err := r.db.Where("id = ?", id).First(&opening).Error
	return &opening, err
}

func (r *careerRepository) UpdateOpening(id string, opening *models.CareerOpening) error {
	return r.db.Model(&models.CareerOpening{}).Where("id = ?", id).Updates(opening).Error
}

func (r *careerRepository) UpdateOpeningStatus(id string, status string) error {
	return r.db.Model(&models.CareerOpening{}).Where("id = ?", id).Update("status", status).Error
}

// ==========================================
// JOB APPLICATION METHODS
// ==========================================

// 👇 This saves the applicant's details and their resume path to the DB
func (r *careerRepository) SubmitApplication(application *models.JobApplication) error {
	return r.db.Create(application).Error
}

func (r *careerRepository) GetAllApplications() ([]models.JobApplication, error) {
	var applications []models.JobApplication
	err := r.db.Preload("CareerOpening").Order("created_at DESC").Find(&applications).Error
	return applications, err
}

func (r *careerRepository) UpdateApplicationStatus(id string, status string) error {
	return r.db.Model(&models.JobApplication{}).Where("id = ?", id).Update("status", status).Error
}