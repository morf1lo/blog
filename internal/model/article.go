package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Article struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"column:user_id" json:"userId"`
	Title     string    `gorm:"column:title" form:"title" json:"title" validate:"required,min=2,max=64"`
	Text      string    `gorm:"column:text" form:"text" json:"text" validate:"required,min=10,max=10000"`
	DateAdded time.Time `gorm:"column:date_added;type:timestamp(0) without time zone;autoCreateTime" json:"dateAdded"`
}

func (a *Article) Validate() error {
	return validator.New().Struct(a)
}

type ArticleDto struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"column:user_id" json:"userId"`
	UserAvatar string    `json:"userAvatar"`
	Username   string    `json:"username"`
	Title      string    `json:"title"`
	Text       string    `json:"text"`
	DateAdded  time.Time `json:"dateAdded"`
}
