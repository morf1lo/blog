package repository

import (
	"github.com/morf1lo/blog/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) FindByUsernameOrEmail(username string, email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ? OR email = ?", username, email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) FindByActivationLink(link string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("activation_link = ?", link).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) FindByToken(token string) (*model.User, error) {
	var user model.User
	if err := r.db.Table("users u").Joins("JOIN reset_password_tokens t ON t.user_email = u.email").Where("t.token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	
	return &user, nil
}

func (r *UserRepo) ExistsByUsernameOrEmail(username string, email string) bool {
	var count int64
	r.db.Model(&model.User{}).Where("username = ? OR email = ?", username, email).Count(&count)
	return count > 0
}

func (r *UserRepo) ExistsByID(id uuid.UUID) bool {
	var count int64
	r.db.Model(&model.User{}).Where("id = ?", id).Count(&count)
	return count > 0
}

func (r *UserRepo) UpdateByID(id uuid.UUID, updates map[string]interface{}) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}
