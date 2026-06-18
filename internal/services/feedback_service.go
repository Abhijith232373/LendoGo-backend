package services

import (
	"errors"

	"lendogo-backend/internal/repositories"
	"lendogo-backend/structures/models"
)

// FeedbackService defines the business logic contract for user feedback operations.
type FeedbackService interface {
	SubmitFeedback(feedback models.Feedback) (*models.Feedback, error)
	GetAllFeedback() ([]models.Feedback, error)
	UpdateStatus(id string, status string) error
}

type feedbackService struct {
	repo repositories.FeedbackRepository
}

// NewFeedbackService injects the repository into the FeedbackService.
func NewFeedbackService(repo repositories.FeedbackRepository) FeedbackService {
	return &feedbackService{
		repo: repo,
	}
}

// SubmitFeedback processes a new feedback entry, applying default business rules.
func (s *feedbackService) SubmitFeedback(feedback models.Feedback) (*models.Feedback, error) {
	// 1. Enforce the default status before saving to the database
	if feedback.Status == "" {
		feedback.Status = "Pending"
	}

	// 2. Pass to the repository
	err := s.repo.CreateFeedback(&feedback)
	if err != nil {
		return nil, errors.New("failed to save feedback, please try again later")
	}

	return &feedback, nil
}

// GetAllFeedback retrieves all feedback entries for the Admin Panel.
func (s *feedbackService) GetAllFeedback() ([]models.Feedback, error) {
	return s.repo.GetAllFeedback()
}

// UpdateStatus validates and updates the status of a specific feedback entry.
func (s *feedbackService) UpdateStatus(id string, status string) error {
	// 1. Business Logic: Prevent invalid statuses from entering the database
	validStatuses := map[string]bool{
		"Pending": true,
		"Replied": true,
		"Ignored": true,
	}
	
	if !validStatuses[status] {
		return errors.New("invalid status: must be 'Pending', 'Replied', or 'Ignored'")
	}

	// 2. Pass the safe data to the repository
	return s.repo.UpdateStatus(id, status)
}