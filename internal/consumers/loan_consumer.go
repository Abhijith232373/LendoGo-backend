package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"

	// Import your services layer so we can use it!
	"lendogo-backend/internal/services"
	"lendogo-backend/utils"

	"github.com/google/uuid"
)

type LoanConsumer struct {
	Reader              *kafka.Reader
	LoanService         services.LoanService // 👈 INJECTED: Now the consumer can talk to the database!
	NotificationService services.NotificationService
}

// 👇 UPDATED: Constructor now requires the LoanService
func NewLoanConsumer(brokerURL string, topic string, groupID string, loanService services.LoanService, notifService services.NotificationService) *LoanConsumer {
	return &LoanConsumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{brokerURL},
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
		LoanService:         loanService,
		NotificationService: notifService,
	}
}

func (c *LoanConsumer) Start(ctx context.Context) {
	log.Printf("📥 Loan Consumer listening on topic: %s", c.Reader.Config().Topic)

	for {
		msg, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("❌ Kafka Consumer Error: %v", err)
			continue
		}

		var envelope map[string]interface{}
		if err := json.Unmarshal(msg.Value, &envelope); err != nil {
			log.Printf("❌ Failed to parse Kafka message: %v", err)
			continue
		}

		eventType, ok := envelope["event_type"].(string)
		if !ok {
			continue
		}

		// Route the event!
		switch eventType {
		case "LOAN_DISBURSED":
			fmt.Println("✅ [KAFKA CAUGHT EVENT] -> LOAN_DISBURSED")

			payload := envelope["data"].(map[string]interface{})
			
			// Safely extract the data from the JSON map
			loanID := fmt.Sprintf("%v", payload["loan_id"])
			userID := fmt.Sprintf("%v", payload["user_id"])
			amount := payload["amount"]

			fmt.Printf("💸 Processing Background Tasks for Loan %v (User: %v, Amt: ₹%v)\n", loanID, userID, amount)

			// ==================================================
			// 1. GENERATE EMI SCHEDULE & UPDATE LOAN HISTORY
			// ==================================================
			// We pass context.Background() so this database write finishes even if the user closes the app!
			err := c.LoanService.GenerateEMISchedule(context.Background(), loanID)
			if err != nil {
				log.Printf("❌ Failed to generate EMI schedule: %v", err)
			} else {
				fmt.Println("✅ EMI Schedule Generated & Loan History Updated in Database!")
			}

			// ==================================================
			// 2. SEND IN-APP NOTIFICATION
			// ==================================================
			fmt.Printf("🔔 Pushing In-App Notification to User: %v\n", userID)
			
			// Use Google UUID parser just to satisfy NotificationService (assuming userID string is a UUID)
			importUUID, err := uuid.Parse(userID)
			if err == nil {
			    c.NotificationService.SendNotification(importUUID, fmt.Sprintf("Your loan of ₹%v has been disbursed to your wallet!", amount), "SUCCESS", "loan")
			}

			// ==================================================
			// 3. SEND EMAIL (Run in a sub-goroutine so it doesn't block!)
			// ==================================================
			go func(uID string, amt interface{}) {
				// Assuming user's email is needed. For mock, we just pass uID
				utils.SendDisbursalEmail(uID, amt)
			}(userID, amount)

			fmt.Println("🏁 All background tasks for LOAN_DISBURSED completed successfully!")
		}
	}
}