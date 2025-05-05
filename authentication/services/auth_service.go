package services

import (
	"authentication/config"
	"authentication/types"
	"authentication/utils"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Signup(user *types.User) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Create user
	result := config.DB.Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *AuthService) Login(email, password string) (*types.User, string, error) {
	var user types.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("user not found")
		}
		return nil, "", err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid password")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return nil, "", err
	}

	// Update last login
	user.LastLogin = time.Now()
	config.DB.Save(&user)

	return &user, token, nil
}

func (s *AuthService) ForgotPassword(email string) error {
	var user types.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	// Generate reset token
	token := utils.GenerateRandomToken()
	user.ResetToken = token
	user.ResetExpires = time.Now().Add(24 * time.Hour)

	if err := config.DB.Save(&user).Error; err != nil {
		return err
	}

	// TODO: Send email with reset token
	return nil
}

func (s *AuthService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	var user types.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return err
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return config.DB.Save(&user).Error
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
	var user types.User
	if err := config.DB.Where("reset_token = ? AND reset_expires > ?", token, time.Now()).First(&user).Error; err != nil {
		return errors.New("invalid or expired reset token")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.ResetToken = ""
	user.ResetExpires = time.Time{}

	return config.DB.Save(&user).Error
}
