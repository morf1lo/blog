package service

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/morf1lo/blog/internal/model"
	"github.com/morf1lo/blog/internal/repository"
	"github.com/morf1lo/blog/pkg/file"
	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type UserService struct {
	repo *repository.Repository
	rdb *redis.Client
}

func NewUserService(repo *repository.Repository, rdb *redis.Client) *UserService {
	return &UserService{
		repo: repo,
		rdb: rdb,
	}
}

func (s *UserService) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.rdb.Get(ctx, userPrefix + id.String()).Result()
	if err != nil {
		if err == redis.Nil {
			userDB, err := s.repo.User.FindByID(id)
			if err != nil {
				return nil, err
			}

			userJSON, err := json.Marshal(userDB)
			if err != nil {
				return nil, err
			}

			if err := s.rdb.Set(ctx, userPrefix + id.String(), userJSON, time.Hour * 2).Err(); err != nil {
				return nil, err
			}

			return userDB, nil
		}

		return nil, err
	}

	var userDB model.User
	if err := json.Unmarshal([]byte(user), &userDB); err != nil {
		return nil, err
	}
	
	return &userDB, nil
}

func (s *UserService) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.rdb.Get(ctx, userNamePrefix + username).Result()
	if err != nil {
		if err == redis.Nil {
			userDB, err := s.repo.User.FindByUsername(username)
			if err != nil {
				return nil, err
			}

			userJSON, err := json.Marshal(userDB)
			if err != nil {
				return nil, err
			}

			if err := s.rdb.Set(ctx, userNamePrefix + username, userJSON, time.Hour * 2).Err(); err != nil {
				return nil, err
			}

			return userDB, nil
		}

		return nil, err
	}

	var userDB model.User
	if err := json.Unmarshal([]byte(user), &userDB); err != nil {
		return nil, err
	}

	return &userDB, nil
}

func (s *UserService) UpdateAvatar(c *gin.Context, userID uuid.UUID, image *multipart.FileHeader) error {
	if !file.IsImage(image) {
		return errFileIsNotAnImage
	}

	userIDString := userID.String()

	uploadPath := "public/avatars/"
	filePath := uploadPath + userIDString + filepath.Ext(image.Filename)

	files, err := filepath.Glob(filepath.Join(uploadPath, userIDString + ".*"))
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	if err := c.SaveUploadedFile(image, filePath); err != nil {
		return err
	}

	avPath := fmt.Sprintf("http://%s:%s/%s", viper.GetString("app.host"), viper.GetString("app.port"), filePath)
	if err := s.repo.User.UpdateByID(userID, map[string]interface{}{"avatar": avPath}); err != nil {
		return err
	}

	if err := s.rdb.Del(c.Request.Context(), userPrefix + userIDString).Err(); err != nil {
		return err
	}

	return nil
}
