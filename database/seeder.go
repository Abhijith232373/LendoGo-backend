package database

import (
	"log"
	"os"

	// Make sure this matches your module name in go.mod!
	"lendogo-backend/structures/models"
)

// RunSeeders is the master function to execute all database seeders
func RunSeeders() {
	log.Println(" Starting database seeders...")

	// 1. Run Admin Seeder
	SeedAdmin()

	// 2. Run System Wallet Seeder
	seedSystemWallet()

	log.Println("✅ All seeders executed successfully.")
}

// SeedAdmin creates the default master admin account if none exists
func SeedAdmin() {
	var adminCount int64

	// 1. Read the credentials from your .env file
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminEmail == "" || adminPassword == "" {
		log.Println("⚠️ WARNING: ADMIN_EMAIL or ADMIN_PASSWORD is not set in .env. Skipping admin creation.")
		return
	}

	// 2. Check if ANY admin already exists in the STAFF table by email
	DB.Model(&models.Staff{}).Where("email = ?", adminEmail).Count(&adminCount)

	if adminCount > 0 {
		// An admin already exists, we don't need to do anything.
		return
	}

	log.Println("No admin found. Creating default master admin account in staffs table...")

	// Note: We DO NOT manually hash the password here using utils.HashPassword.
	// Because your models.Staff has a BeforeCreate hook, GORM will automatically
	// hash this plain text password right before saving it to the database!
	
	// 3. Build the Admin Staff Member
	adminUser := models.Staff{
		FullName: "System Administrator",
		Email:    adminEmail,
		Password: adminPassword, // Passed as plain text!
		Role:     "Superadmin",  // Triggers the Master Bypass in our Middleware
		Status:   "Active",
		Permissions: map[string]bool{
			"dashboard.view":      true,
			"users.read":          true,
			"users.create":        true,
			"users.update":        true,
			"users.delete":        true,
			"loans.view":          true,
			"loans.update":        true,
			"consultation.view":   true,
		},
	}

	// 4. Save to the database
	if err := DB.Create(&adminUser).Error; err != nil {
		log.Printf("⚠️ Warning: Failed to seed admin user: %v\n", err)
	} else {
		log.Println("✅ Default Admin account created successfully in staffs table!")
	}
}

// seedSystemWallet ensures the Master Capital Ledger exists for Admin disbursements
func seedSystemWallet() {
	var wallet models.SystemWallet

	// FirstOrCreate checks if a row with WalletName="capital_disbursement" exists.
	// If not, it creates it with a 0.0 balance.
	result := DB.FirstOrCreate(&wallet, models.SystemWallet{
		WalletName: "capital_disbursement",
		Balance:    0.0,
	})

	if result.Error != nil {
		log.Printf("❌ Failed to seed System Wallet: %v\n", result.Error)
	} else {
		log.Println("💰 System Wallet verified/seeded.")
	}
}