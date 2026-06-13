package services

import (
	"encoding/json" 
	"errors"

	"lendogo-backend/internal/repositories"
	"lendogo-backend/structures/models"
)

type CareerService interface {
	// Job Opening Methods
	CreateOpening(req models.CareerOpening) (*models.CareerOpening, error)
	GetAllOpenings(status string) ([]models.CareerOpening, error)
	GetOpeningByID(id string) (*models.CareerOpening, error)
	UpdateOpening(id string, req models.CareerOpening) (*models.CareerOpening, error)
	UpdateOpeningStatus(id string, status string) error
	
	// 👇 The new Application Method
	SubmitApplication(req models.JobApplication) (*models.JobApplication, error)
	GetAllApplications() ([]models.JobApplication, error)
	UpdateApplicationStatus(id string, status string) error
}

type careerService struct {
	repo repositories.CareerRepository
}

func NewCareerService(repo repositories.CareerRepository) CareerService {
	return &careerService{repo: repo}
}

// ==========================================
// JOB OPENING LOGIC
// ==========================================

func (s *careerService) CreateOpening(req models.CareerOpening) (*models.CareerOpening, error) {
	validModes := map[string]bool{"Hybrid": true, "Remote": true, "On-site": true}
	if !validModes[req.WorkMode] {
		return nil, errors.New("invalid work mode: must be Hybrid, Remote, or On-site")
	}

	if req.Skills == nil { req.Skills = json.RawMessage("[]") }
	if req.Responsibilities == nil { req.Responsibilities = json.RawMessage("[]") }
	if req.Requirements == nil { req.Requirements = json.RawMessage("[]") }
	if req.Benefits == nil { req.Benefits = json.RawMessage("[]") }

	if req.Status == "" {
		req.Status = "Open"
	}

	err := s.repo.CreateOpening(&req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *careerService) GetAllOpenings(status string) ([]models.CareerOpening, error) {
	return s.repo.GetAllOpenings(status)
}

func (s *careerService) GetOpeningByID(id string) (*models.CareerOpening, error) {
	return s.repo.GetOpeningByID(id)
}

func (s *careerService) UpdateOpening(id string, req models.CareerOpening) (*models.CareerOpening, error) {
	if req.Skills == nil { req.Skills = json.RawMessage("[]") }
	if req.Responsibilities == nil { req.Responsibilities = json.RawMessage("[]") }
	if req.Requirements == nil { req.Requirements = json.RawMessage("[]") }
	if req.Benefits == nil { req.Benefits = json.RawMessage("[]") }

	err := s.repo.UpdateOpening(id, &req)
	if err != nil {
		return nil, err
	}
	return s.GetOpeningByID(id)
}

func (s *careerService) UpdateOpeningStatus(id string, status string) error {
	return s.repo.UpdateOpeningStatus(id, status)
}

// ==========================================
// JOB APPLICATION LOGIC
// ==========================================

// 👇 SubmitApplication verifies the job exists before saving the candidate
func (s *careerService) SubmitApplication(req models.JobApplication) (*models.JobApplication, error) {
	// 1. Verify the job actually exists before allowing an application
	_, err := s.repo.GetOpeningByID(req.CareerOpeningID.String())
	if err != nil {
		return nil, errors.New("the job opening you are applying for does not exist or was closed")
	}

	// 2. Default status for HR
	req.Status = "Reviewing"

	// 3. Save to database
	if err := s.repo.SubmitApplication(&req); err != nil {
		return nil, errors.New("failed to submit application, please try again")
	}

	return &req, nil
}

func (s *careerService) GetAllApplications() ([]models.JobApplication, error) {
	return s.repo.GetAllApplications()
}

func (s *careerService) UpdateApplicationStatus(id string, status string) error {
	return s.repo.UpdateApplicationStatus(id, status)
}