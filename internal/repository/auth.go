package repository

import (
	"github.com/morf1lo/blog/internal/model"
	"gorm.io/gorm"
)

type AuthRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (s *AuthRepo) CreateUser(user *model.User) error {
	return s.db.Save(user).Error
}

func (r *AuthRepo) SaveResetPasswordToken(token *model.ResetPasswordToken) error {
	return r.db.Save(token).Error
}

func (r *AuthRepo) FindResetPasswordToken(token string) (*model.ResetPasswordToken, error) {
	var tokenDB model.ResetPasswordToken
	if err := r.db.Where("token = ?", token).First(&tokenDB).Error; err != nil {
		return nil, err
	}

	return &tokenDB, nil
}

func (r *AuthRepo) DeleteResetPasswordToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&model.ResetPasswordToken{}).Error
}
