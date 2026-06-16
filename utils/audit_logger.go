package utils

import (
	"log"
	"lendogo-backend/database"
	"lendogo-backend/internal/websockets"
	"lendogo-backend/structures/models"
	"github.com/google/uuid"
)

func RecordAudit(actorID uuid.UUID, actorName, actionType, entityType, entityID, description, ipAddress string) {
	auditEntry := models.AuditLog{
		ActorID:     actorID,
		ActorName:   actorName,
		ActionType:  actionType,
		EntityType:  entityType,
		EntityID:    entityID,
		Description: description,
		IPAddress:   ipAddress,
	}
	if err := database.DB.Create(&auditEntry).Error; err != nil {
		log.Printf("CRITICAL COMPLIANCE ERROR: Failed to write audit log: %v\n", err)
		return 
	}
	websockets.BroadcastMessage("NEW_AUDIT_LOG", auditEntry)
}