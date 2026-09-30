package main

import (
	"fmt"
	"log"

	"golang-test/config"
	"golang-test/database"
	"golang-test/internal/handler"
	"golang-test/internal/repository"
	"golang-test/internal/service"
	"golang-test/pkg/oauth"
	"golang-test/routes"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Initialize Database Connection
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 3. Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	eventRepo := repository.NewEventRepository(db)

	// 4. Initialize OAuth Utility
	googleOAuth := oauth.NewGoogleOAuth(cfg)

	// 5. Initialize Services
	authService := service.NewAuthService(userRepo, refreshTokenRepo, cfg, googleOAuth)
	calendarService := service.NewCalendarService(eventRepo, userRepo, googleOAuth)

	// 6. Initialize Handlers
	authHandler := handler.NewAuthHandler(authService)
	calendarHandler := handler.NewCalendarHandler(calendarService)

	// 7. Setup Router
	r := routes.SetupRouter(cfg, authHandler, calendarHandler)

	// 8. Start HTTP Server
	serverAddr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server running on port %s in %s mode...", cfg.AppPort, cfg.AppEnv)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
