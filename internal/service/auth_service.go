package service

import (
	"context"
	"errors"
	"time"

	"golang-test/config"
	"golang-test/internal/dto"
	"golang-test/internal/model"
	"golang-test/internal/repository"
	"golang-test/pkg/auth"
	"golang-test/pkg/oauth"
)

var (
	ErrEmailAlreadyExists = errors.New("email is already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid or expired refresh token")
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(req dto.RefreshTokenRequest) (*dto.AuthResponse, error)
	Logout(refreshToken string) error
	GetProfile(userID int64) (*dto.UserResponse, error)
	GetGoogleAuthURL() string
	HandleGoogleCallback(ctx context.Context, code string) (*dto.AuthResponse, error)
}

type authService struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	cfg              *config.Config
	googleOAuth      *oauth.GoogleOAuth
}

func NewAuthService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	cfg *config.Config,
	googleOAuth *oauth.GoogleOAuth,
) AuthService {
	return &authService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		cfg:              cfg,
		googleOAuth:      googleOAuth,
	}
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: &hashedPassword,
		Provider:     "email",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return s.generateTokensAndResponse(user)
}

func (s *authService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil || user.PasswordHash == nil {
		return nil, ErrInvalidCredentials
	}

	if !auth.CheckPasswordHash(req.Password, *user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	return s.generateTokensAndResponse(user)
}

func (s *authService) RefreshToken(req dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claims, err := auth.ValidateToken(req.RefreshToken, s.cfg.JWTSecret)
	if err != nil {
		return nil, ErrInvalidToken
	}

	tokenRecord, err := s.refreshTokenRepo.FindByToken(req.RefreshToken)
	if err != nil || tokenRecord == nil {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	_ = s.refreshTokenRepo.DeleteByToken(req.RefreshToken)

	return s.generateTokensAndResponse(user)
}

func (s *authService) Logout(refreshToken string) error {
	return s.refreshTokenRepo.DeleteByToken(refreshToken)
}

func (s *authService) GetProfile(userID int64) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	userResp := s.toUserResponse(user)
	return &userResp, nil
}

func (s *authService) GetGoogleAuthURL() string {
	return s.googleOAuth.GetAuthURL("state-token")
}

func (s *authService) HandleGoogleCallback(ctx context.Context, code string) (*dto.AuthResponse, error) {
	token, err := s.googleOAuth.ExchangeCode(ctx, code)
	if err != nil {
		return nil, err
	}

	googleUser, err := s.googleOAuth.GetGoogleUserInfo(ctx, token)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByEmail(googleUser.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		providerID := googleUser.ID
		avatar := googleUser.Picture
		user = &model.User{
			Name:       googleUser.Name,
			Email:      googleUser.Email,
			Provider:   "google",
			ProviderID: &providerID,
			Avatar:     &avatar,
		}
		if err := s.userRepo.Create(user); err != nil {
			return nil, err
		}
	} else {
		providerID := googleUser.ID
		avatar := googleUser.Picture
		user.ProviderID = &providerID
		user.Avatar = &avatar
		_ = s.userRepo.Update(user)
	}

	_ = s.userRepo.UpdateGoogleTokens(user.ID, token.AccessToken, token.RefreshToken, token.Expiry)

	return s.generateTokensAndResponse(user)
}

func (s *authService) generateTokensAndResponse(user *model.User) (*dto.AuthResponse, error) {
	accessToken, err := auth.GenerateAccessToken(user.ID, user.Email, s.cfg.JWTSecret, s.cfg.JWTAccessExpirationHours)
	if err != nil {
		return nil, err
	}

	refreshToken, expiresAt, err := auth.GenerateRefreshToken(user.ID, user.Email, s.cfg.JWTSecret, s.cfg.JWTRefreshExpirationDays)
	if err != nil {
		return nil, err
	}

	tokenRecord := &model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	}
	if err := s.refreshTokenRepo.Create(tokenRecord); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         s.toUserResponse(user),
	}, nil
}

func (s *authService) toUserResponse(user *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Provider:  user.Provider,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
}
