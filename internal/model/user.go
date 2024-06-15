package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `gorm:"primaryKey" json:"id"`
	Username       string    `gorm:"column:username;unique" form:"username" json:"username" validate:"required,min=3,max=20"`
	Email          string    `gorm:"column:email;unique" form:"email" json:"email" validate:"required,email"`
	Password       string    `gorm:"column:password" form:"password" json:"password" validate:"required,min=8,max=24"`
	Avatar         string    `gorm:"column:avatar" json:"avatar"`
	Activated      bool      `gorm:"column:activated" json:"activated"`
	ActivationLink *string   `gorm:"column:activation_link;unique" json:"activationLink"`
	DateAdded      time.Time `gorm:"column:date_added;type:timestamp(0) without time zone;autoCreateTime" json:"dateAdded"`
}

func (u *User) Validate() error {
	return validator.New().Struct(u)
}

func (u *User) DTO() *User {
	return &User{
		ID: u.ID,
		Username: u.Username,
		Email: u.Email,
		// IGNORING u.Password
		Avatar: u.Avatar,
		Activated: u.Activated,
		// IGNORING u.ActivationLink
		DateAdded: u.DateAdded,
	}
}
