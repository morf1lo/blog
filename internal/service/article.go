package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/morf1lo/blog/internal/model"
	"github.com/morf1lo/blog/internal/repository"
	"github.com/redis/go-redis/v9"
)

type ArticleService struct {
	repo *repository.Repository
	rdb *redis.Client
}

func NewArticleRepo(repo *repository.Repository, rdb *redis.Client) *ArticleService {
	return &ArticleService{
		repo: repo,
		rdb: rdb,
	}
}

const (
	articleExpiry = time.Minute * 10
	lastArticlesExpiry = time.Minute * 2
)

func (s *ArticleService) Create(ctx context.Context, article *model.Article) error {
	if err := article.Validate(); err != nil {
		return err
	}

	article.Title = strings.TrimSpace(article.Title)
	article.Text = strings.TrimSpace(article.Text)

	if err := s.rdb.Del(ctx, authorArticlesPrefix + article.UserID.String()).Err(); err != nil {
		return err
	}

	return s.repo.Article.Create(article)
}

func (s *ArticleService) FindByID(ctx context.Context, id int64) (*model.ArticleDto, error) {
	idString := strconv.Itoa(int(id))

	article, err := s.rdb.Get(ctx, articlePrefix + idString).Result()
	if err != nil {
		if err == redis.Nil {
			articleDB, err := s.repo.Article.FindByID(id)
			if err != nil {
				return nil, err
			}

			articleJSON, err := json.Marshal(articleDB)
			if err != nil {
				return nil, err
			}

			if err := s.rdb.Set(ctx, articlePrefix + idString, articleJSON, articleExpiry).Err(); err != nil {
				return nil, err
			}

			return articleDB, nil
		}

		return nil, err
	}

	var articleDB model.ArticleDto
	if err := json.Unmarshal([]byte(article), &articleDB); err != nil {
		return nil, err
	}

	return &articleDB, nil
}

func (s *ArticleService) Search(ctx context.Context, query string) ([]*model.Article, error) {
	queryLow := strings.ToLower(query)

	articles, err := s.rdb.Get(ctx, articlesSearchResultsFor + queryLow).Result()
	if err != nil {
		if err == redis.Nil {
			articlesDB, err := s.repo.Article.Search(queryLow)
			if err != nil {
				return nil, err
			}

			articlesJSON, err := json.Marshal(articlesDB)
			if err != nil {
				return nil, err
			}

			if err := s.rdb.Set(ctx, articlesSearchResultsFor + queryLow, articlesJSON, time.Minute * 10).Err(); err != nil {
				return nil, err
			}

			return articlesDB, nil
		}

		return nil, err
	}

	var articlesDB []*model.Article
	if err := json.Unmarshal([]byte(articles), &articlesDB); err != nil {
		return nil, err
	}

	return articlesDB, nil
}

func (s *ArticleService) FindAuthorArticles(ctx context.Context, authorID string) ([]*model.Article, error) {
	articles, err := s.rdb.Get(ctx, authorArticlesPrefix + authorID).Result()
	if err != nil {
		if err == redis.Nil {
			articlesDB, err := s.repo.Article.FindAuthorArticles(authorID)
			if err != nil {
				return nil, err
			}

			articlesJSON, err := json.Marshal(articlesDB)
			if err != nil {
				return nil, err
			}

			if err := s.rdb.Set(ctx, authorArticlesPrefix + authorID, articlesJSON, time.Hour * 1).Err(); err != nil {
				return nil, err
			}

			return articlesDB, nil
		}

		return nil, err
	}

	var articlesDB []*model.Article
	if err := json.Unmarshal([]byte(articles), &articlesDB); err != nil {
		return nil, err
	}

	return articlesDB, nil
}

func (s *ArticleService) FindLastArticles(ctx context.Context, limit int) ([]*model.ArticleDto, error) {
	articles, err := s.rdb.Get(ctx, lastArticlesPrefix).Result()
	if err != nil {
		if err == redis.Nil {
			articlesDB, err := s.repo.Article.FindLastArticles(limit)
			if err != nil {
				return nil, err
			}

			articlesJSON, err := json.Marshal(articlesDB)
			if err != nil {
				return nil, err
			}

			if err := s.rdb.Set(ctx, lastArticlesPrefix, articlesJSON, lastArticlesExpiry).Err(); err != nil {
				return nil, err
			}

			return articlesDB, nil
		}

		return nil, err
	}

	var articlesDB []*model.ArticleDto
	if err := json.Unmarshal([]byte(articles), &articlesDB); err != nil {
		return nil, err
	}

	return articlesDB, nil
}
