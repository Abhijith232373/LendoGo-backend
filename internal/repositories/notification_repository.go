package repositories

import (
	"lendogo-backend/structures/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	CreateNotification(notification *models.Notification) error
	GetUnreadNotifications(userID uuid.UUID) ([]models.Notification, error)
	MarkAsRead(userID uuid.UUID) error
}

type notificationRepoImpl struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepoImpl{db: db}
}

func (r *notificationRepoImpl) CreateNotification(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepoImpl) GetUnreadNotifications(userID uuid.UUID) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("user_id = ? AND is_read = ?", userID, false).Order("created_at desc").Limit(20).Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepoImpl) MarkAsRead(userID uuid.UUID) error {
	return r.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}
