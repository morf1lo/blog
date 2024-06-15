package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Comment struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"column:user_id" json:"userId"`
	ArticleID int64     `gorm:"column:article_id" form:"articleId" json:"articleId" validate:"required"`
	Body      string    `gorm:"column:body" form:"body" json:"body" validate:"required,min=1,max=512"`
	DateAdded time.Time `gorm:"column:date_added;type:timestamp(0) without time zone;autoCreateTime" json:"dateAdded"`
}

func (c *Comment) Validate() error {
	return validator.New().Struct(c)
}

type CommentDto struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"column:user_id" json:"userId"`
	UserAvatar string    `json:"userAvatar"`
	Username   string    `json:"username"`
	ArticleID  int64     `json:"articleId"`
	Body       string    `json:"body"`
	DateAdded  time.Time `json:"dateAdded"`
}
