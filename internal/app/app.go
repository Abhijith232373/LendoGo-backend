package app

import (
	"context" // 👈 Added for Kafka background context
	"os"      // 👈 Added to read the .env KAFKA_BROKER

	"github.com/gofiber/fiber/v2"

	"lendogo-backend/database"
	"lendogo-backend/internal/controllers/admin_controller"
	"lendogo-backend/internal/controllers/auth_controller"
	"lendogo-backend/internal/controllers/chat_controller"
	consultation_controller "lendogo-backend/internal/controllers/consultation_controller"
	"lendogo-backend/internal/controllers/loan_controller"
	"lendogo-backend/internal/controllers/user_profile_controller"
	"lendogo-backend/internal/controllers/wallet_controller"
	"lendogo-backend/internal/repositories"
	"lendogo-backend/internal/routes"
	"lendogo-backend/internal/services"

	"lendogo-backend/internal/consumers"
	"lendogo-backend/utils" // 👇 THIS IS NOW UNCOMMENTED
)

// SetupApp initializes all dependencies and registers routes
func SetupApp(app *fiber.App) {

	// ==========================================
	// 1. KAFKA INFRASTRUCTURE SETUP 🚀
	// ==========================================
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	// 👇 THIS IS NOW UNCOMMENTED
	kafkaProducer := utils.NewKafkaProducer(kafkaBroker)

	// ==========================================
	// 2. REPOSITORIES (Data Layer)
	// ==========================================
	userRepo := repositories.NewUserRepository(database.DB)
	consultationRepo := repositories.NewConsultationRepository(database.DB)
	loanRepo := repositories.NewLoanRepository(database.DB)
	walletRepo := repositories.NewWalletRepository(database.DB)
	chatRepo := repositories.NewChatRepository(database.DB)
	profileRepo := repositories.NewUserProfileRepository(database.DB)

	// ==========================================
	// 3. SERVICES & HUBS (Business Logic Layer)
	// ==========================================
	authService := services.NewAuthService(userRepo)
	consultationService := services.NewConsultationService(consultationRepo)
	loanService := services.NewLoanService(loanRepo)

	// 👇 THIS LINE IS THE FIX: We are now passing the kafkaProducer into the service!
	walletService := services.NewWalletService(walletRepo, kafkaProducer)

	profileService := services.NewUserProfileService(profileRepo)

	chatHub := services.NewChatHub(chatRepo)
	go chatHub.Run()

	// ==========================================
	// 4. KAFKA CONSUMERS (Background Workers) 📥
	// ==========================================
	paymentConsumer := consumers.NewPaymentConsumer(
		kafkaBroker,
		"telemetry.payments",
		"payment-processor-group",
		loanService,
	)

	// 🔥 Start the consumer safely in the background!
	go paymentConsumer.Start(context.Background())

	// ==========================================
	// 5. CONTROLLERS (HTTP Layer)
	// ==========================================
	authController := auth_controller.NewAuthController(authService)
	consultationController := consultation_controller.NewConsultationController(consultationService)
	adminController := admin_controller.NewAdminController()
	loanController := loan_controller.NewLoanController(loanService)
	walletController := wallet_controller.NewWalletController(walletService)
	chatController := chat_controller.NewChatController(chatHub)
	profileController := user_profile_controller.NewUserProfileController(profileService)

	// ==========================================
	// 6. ROUTER SETUP
	// ==========================================
	api := app.Group("/api")

	routes.SetupAuthRoutes(api, authController)
	routes.SetupConsultationRoutes(api, consultationController)
	routes.SetupAdminRoutes(api, adminController)
	routes.SetupLoanRoutes(api, loanController)
	routes.SetupWalletRoutes(api, walletController)
	routes.SetupChatRoutes(api, chatController)
	routes.SetupUserProfileRoutes(api, profileController)
}