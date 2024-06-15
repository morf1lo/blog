package model

import "time"

type ResetPasswordToken struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserEmail string    `gorm:"column:user_email;unique" json:"userEmail" binding:"required,email"`
	Token     string    `gorm:"column:token;unique" json:"token"`
	Expiry    time.Time `gorm:"column:expiry" json:"expiry"`
}
