package repositories

import (
	"errors"
	"lendogo-backend/structures/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserProfileRepository interface {
	GetProfile(userID uuid.UUID) (*models.UserProfile, *models.User, error)
	UpsertProfile(userID uuid.UUID, profile models.UserProfile, newFullName string) error
}

type userProfileRepoImpl struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) UserProfileRepository {
	return &userProfileRepoImpl{db: db}
}

func (r *userProfileRepoImpl) GetProfile(userID uuid.UUID) (*models.UserProfile, *models.User, error) {
	var profile models.UserProfile
	var user models.User

	// Get Auth User details (for full name and email)
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, nil, err
	}

	// Get Profile details (If not found, we just return empty profile without error)
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, err
	}

	return &profile, &user, nil
}

func (r *userProfileRepoImpl) UpsertProfile(userID uuid.UUID, profile models.UserProfile, newFullName string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Update the full_name in the main users table
		if newFullName != "" {
			if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("full_name", newFullName).Error; err != nil {
				return err
			}
		}

		// 2. Check if profile exists
		var existing models.UserProfile
		err := tx.Where("user_id = ?", userID).First(&existing).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new
			profile.UserID = userID
			return tx.Create(&profile).Error
		} else if err != nil {
			return err
		}

		// 3. Update existing profile (keeping old image if a new one wasn't uploaded)
		updates := map[string]interface{}{
			"phone_number":  profile.PhoneNumber,
			"date_of_birth": profile.DateOfBirth,
			"pincode":       profile.Pincode,
			"address":       profile.Address,
		}
		if profile.ProfileImage != "" {
			updates["profile_image"] = profile.ProfileImage
		}

		return tx.Model(&existing).Updates(updates).Error
	})
}