package repository

import (
	"github.com/morf1lo/blog/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Authorization interface {
	CreateUser(user *model.User) error
	SaveResetPasswordToken(token *model.ResetPasswordToken) error
	FindResetPasswordToken(token string) (*model.ResetPasswordToken, error)
	DeleteResetPasswordToken(token string) error
}

type User interface {
	FindByID(id uuid.UUID) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByUsernameOrEmail(username string, email string) (*model.User, error)
	FindByActivationLink(link string) (*model.User, error)
	FindByToken(token string) (*model.User, error)
	ExistsByUsernameOrEmail(username string, email string) bool
	ExistsByID(id uuid.UUID) bool
	UpdateByID(id uuid.UUID, updates map[string]interface{}) error
}

type Article interface {
	Create(article *model.Article) error
	FindByID(id int64) (*model.ArticleDto, error)
	UpdateByID(id int64, updates map[string]interface{}) error
	Search(query string) ([]*model.Article, error)
	FindAuthorArticles(authorID string) ([]*model.Article, error)
	FindLastArticles(limit int) ([]*model.ArticleDto, error)
	Delete(id int64, userID uuid.UUID) error
}

type Comment interface {
	Create(comment *model.Comment) error
	FindByID(id int64) (*model.Comment, error)
	FindArticleComments(articleID int64) ([]*model.CommentDto, error)
	Delete(commentID int64) error
}

type Repository struct {
	Authorization
	User
	Article
	Comment
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepo(db),
		User: NewUserRepo(db),
		Article: NewArticleRepo(db),
		Comment: NewCommentRepo(db),
	}
}
