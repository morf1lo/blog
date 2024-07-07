package service

import (
	"context"
	"mime/multipart"

	"github.com/morf1lo/blog/internal/model"
	"github.com/morf1lo/blog/internal/mq"
	"github.com/morf1lo/blog/internal/repository"
	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Authorization interface {
	SignUp(user *model.User) (string, error)
	SignIn(user *model.User) (string, error)
	Activate(ctx context.Context, activationLink string) error
	SaveResetPasswordToken(token *model.ResetPasswordToken) error
	ResetPassword(token string, newPassword string) error
}

type Mail interface {
	ProcessActivationMails()
	SendActivationMail(to []string, link string) error
	ProcessResetPasswordTokenMails()
	SendResetPasswordToken(to []string, link string) error
}

type User interface {
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	UpdateAvatar(c *gin.Context, userID uuid.UUID, file *multipart.FileHeader) error
}

type Article interface {
	Create(ctx context.Context, article *model.Article) error
	FindByID(ctx context.Context, id int64) (*model.ArticleDto, error)
	Search(ctx context.Context, query string) ([]*model.Article, error)
	FindAuthorArticles(ctx context.Context, authorID string) ([]*model.Article, error)
	FindLastArticles(ctx context.Context, limit int) ([]*model.ArticleDto, error)
}

type Comment interface {
	Create(ctx context.Context, comment *model.Comment) error
	FindArticleComments(ctx context.Context, articleID int64) ([]*model.CommentDto, error)
	Delete(ctx context.Context, commentID int64, userID uuid.UUID) error
}

type Service struct {
	Authorization
	Mail
	User
	Article
	Comment
}

func New(repo *repository.Repository, rdb *redis.Client, rabbitMQ *mq.MQConn) *Service {
	return &Service{
		Authorization: NewAuthService(repo, rabbitMQ, rdb),
		Mail: NewMailService(rabbitMQ),
		User: NewUserService(repo, rdb),
		Article: NewArticleRepo(repo, rdb),
		Comment: NewCommentService(repo, rdb),
	}
}
