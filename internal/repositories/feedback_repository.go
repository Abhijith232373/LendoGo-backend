package repositories

import (
	"gorm.io/gorm"

	"lendogo-backend/structures/models"
)

// FeedbackRepository defines the data access contract for user feedback operations.
type FeedbackRepository interface {
	CreateFeedback(feedback *models.Feedback) error
	GetAllFeedback() ([]models.Feedback, error)
	UpdateStatus(id string, status string) error
}

type feedbackRepository struct {
	db *gorm.DB
}

// NewFeedbackRepository allocates and returns a new FeedbackRepository instance.
func NewFeedbackRepository(db *gorm.DB) FeedbackRepository {
	return &feedbackRepository{
		db: db,
	}
}

// CreateFeedback persists a new feedback record to the database.
func (r *feedbackRepository) CreateFeedback(feedback *models.Feedback) error {
	return r.db.Create(feedback).Error
}

// GetAllFeedback retrieves all feedback records ordered by creation date descending.
// Preloads the User entity to strictly prevent N+1 query inefficiencies during admin panel rendering.
func (r *feedbackRepository) GetAllFeedback() ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	err := r.db.Preload("User").Order("created_at DESC").Find(&feedbacks).Error
	if err != nil {
		return nil, err
	}
	return feedbacks, nil
}

// UpdateStatus mutates the status column of a specific feedback record via its ID.
// Uses a targeted UPDATE statement to avoid full-record memory allocation.
func (r *feedbackRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&models.Feedback{}).Where("id = ?", id).Update("status", status).Error
}