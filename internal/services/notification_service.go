package services

import (
	"lendogo-backend/internal/repositories"
	"lendogo-backend/structures/models"

	"github.com/google/uuid"
)

type NotificationService interface {
	SendNotification(userID uuid.UUID, message string, notifType string, target string) error
	GetUnreadNotifications(userID uuid.UUID) ([]models.Notification, error)
	MarkAsRead(userID uuid.UUID) error
}

type notificationServiceImpl struct {
	repo repositories.NotificationRepository
}

func NewNotificationService(repo repositories.NotificationRepository) NotificationService {
	return &notificationServiceImpl{repo: repo}
}

func (s *notificationServiceImpl) SendNotification(userID uuid.UUID, message string, notifType string, target string) error {
	notif := &models.Notification{
		UserID:  userID,
		Message: message,
		Type:    notifType,
		Target:  target,
	}
	return s.repo.CreateNotification(notif)
}

func (s *notificationServiceImpl) GetUnreadNotifications(userID uuid.UUID) ([]models.Notification, error) {
	return s.repo.GetUnreadNotifications(userID)
}

func (s *notificationServiceImpl) MarkAsRead(userID uuid.UUID) error {
	return s.repo.MarkAsRead(userID)
}
