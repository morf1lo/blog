package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/morf1lo/blog/internal/model"
	"github.com/morf1lo/blog/internal/repository"
	"github.com/redis/go-redis/v9"

	"github.com/google/uuid"
)

type CommentService struct {
	repo *repository.Repository
	rdb *redis.Client
}

func NewCommentService(repo *repository.Repository, rdb *redis.Client) *CommentService {
	return &CommentService{
		repo: repo,
		rdb: rdb,
	}
}

const articleCommentsExpiry = time.Minute * 10

func (s *CommentService) Create(ctx context.Context, comment *model.Comment) error {
	if err := comment.Validate(); err != nil {
		return err
	}

	articleIDString := strconv.Itoa(int(comment.ArticleID))

	if err := s.rdb.Del(ctx, articleCommentsPrefix + articleIDString).Err(); err != nil {
		return err
	}

	return s.repo.Comment.Create(comment)
}

func (s *CommentService) FindArticleComments(ctx context.Context, articleID int64) ([]*model.CommentDto, error) {
	articleIDString := strconv.Itoa(int(articleID))

	comments, err := s.rdb.Get(ctx, articleCommentsPrefix + articleIDString).Result()
	if err != nil {
		if err == redis.Nil {
			commentsDB, err := s.repo.Comment.FindArticleComments(articleID)
			if err != nil {
				return nil, err
			}

			commentsJSON, err := json.Marshal(commentsDB)
			if err != nil {
				return nil, err
			}

			if err := s.rdb.Set(ctx, articleCommentsPrefix + articleIDString, commentsJSON, articleCommentsExpiry).Err(); err != nil {
				return nil, err
			}

			return commentsDB, nil
		}

		return nil, err
	}

	var commentsDB []*model.CommentDto
	if err := json.Unmarshal([]byte(comments), &commentsDB); err != nil {
		return nil, err
	}

	return commentsDB, nil
}

func (s *CommentService) Delete(ctx context.Context, commentID int64, userID uuid.UUID) error {
	comment, err := s.repo.Comment.FindByID(commentID)
	if err != nil {
		return err
	}

	article, err := s.repo.Article.FindByID(comment.ArticleID)
	if err != nil {
		return err
	}

	if userID == article.UserID || userID == comment.UserID {
		articleIDString := strconv.Itoa(int(article.ID))

		if err := s.rdb.Del(ctx, articleCommentsPrefix + articleIDString).Err(); err != nil {
			return err
		}

		return s.repo.Comment.Delete(commentID)
	}
	
	return errNoAccess
}
