package jobs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	
	// Import your services so the job can talk to the database!
	"lendogo-backend/internal/services"
)

type EMICheckerJob struct {
	cron        *cron.Cron
	LoanService services.LoanService
}

func NewEMICheckerJob(loanService services.LoanService) *EMICheckerJob {
	// Create a new cron instance based on the server's local time
	c := cron.New(cron.WithLocation(time.Local))
	return &EMICheckerJob{
		cron:        c,
		LoanService: loanService,
	}
}

func (j *EMICheckerJob) Start() {
	// ⏰ THE CRON EXPRESSION: "0 0 * * *" means run exactly at Midnight (12:00 AM) every day.
	// For testing purposes, we can change this to "* * * * *" to run every 1 minute!
	_, err := j.cron.AddFunc("0 0 * * *", func() {
		fmt.Println("⏰ [CRON WAKING UP] Checking database for due EMIs...")

		// 1. Fetch all pending EMIs due today and update to OVERDUE
		// 2. Send "Payment Due" or "Overdue" notifications
		
		// Create a context and call the service layer
		ctx := context.Background()
		err := j.LoanService.ProcessDueEMIs(ctx)
		if err != nil {
			log.Printf("❌ EMI Checker encountered an error: %v", err)
		}

		fmt.Println("🏁 [CRON FINISHED] All EMI reminders processed successfully!")
	})

	if err != nil {
		log.Fatalf("❌ Failed to start EMI Checker Cron Job: %v", err)
	}

	j.cron.Start()
	log.Println("⏱️ EMI Checker Job Scheduled! (Runs daily at Midnight)")
}