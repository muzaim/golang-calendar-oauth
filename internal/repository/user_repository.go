package repository

import (
	"errors"
	"time"

	"golang-test/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id int64) (*model.User, error)
	FindByProviderID(provider string, providerID string) (*model.User, error)
	UpdateGoogleTokens(userID int64, accessToken, refreshToken string, expiresAt time.Time) error
	Update(user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id int64) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByProviderID(provider string, providerID string) (*model.User, error) {
	var user model.User
	err := r.db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateGoogleTokens(userID int64, accessToken, refreshToken string, expiresAt time.Time) error {
	updates := map[string]interface{}{
		"google_access_token":     accessToken,
		"google_token_expires_at": expiresAt,
	}
	if refreshToken != "" {
		updates["google_refresh_token"] = refreshToken
	}
	return r.db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}
