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
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	eventRepo := repository.NewEventRepository(db)

	googleOAuth := oauth.NewGoogleOAuth(cfg)

	authService := service.NewAuthService(userRepo, refreshTokenRepo, cfg, googleOAuth)
	calendarService := service.NewCalendarService(eventRepo, userRepo, googleOAuth)

	authHandler := handler.NewAuthHandler(authService)
	calendarHandler := handler.NewCalendarHandler(calendarService)

	r := routes.SetupRouter(cfg, authHandler, calendarHandler)

	serverAddr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server running on port %s in %s mode...", cfg.AppPort, cfg.AppEnv)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
