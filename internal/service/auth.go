package service

import (
	"context"
	"encoding/json"

	"github.com/morf1lo/blog/internal/model"
	"github.com/morf1lo/blog/internal/mq"
	"github.com/morf1lo/blog/internal/repository"
	"github.com/morf1lo/blog/pkg/auth"
	"github.com/morf1lo/blog/pkg/exp"
	"github.com/morf1lo/blog/pkg/hasher"
	"github.com/redis/go-redis/v9"

	"github.com/google/uuid"
)

type AuthService struct {
	repo *repository.Repository
	rdb *redis.Client
	rabbitMQ *mq.MQConn
}

func NewAuthService(repo *repository.Repository, rabbitMQ *mq.MQConn, rdb *redis.Client) *AuthService {
	return &AuthService{
		repo: repo,
		rdb: rdb,
		rabbitMQ: rabbitMQ,
	}
}

func (s *AuthService) SignUp(user *model.User) (string, error) {
	if err := user.Validate(); err != nil {
		return "", err
	}

	if exists := s.repo.User.ExistsByUsernameOrEmail(user.Username, user.Email); exists {
		return "", errUserIsAlreadyExists
	}

	user.ID = uuid.New()

	activationLink, err := hasher.NewHash(15)
	if err != nil {
		return "", err
	}
	user.ActivationLink = &activationLink

	passwordHash, err := hasher.HashPassword([]byte(user.Password))
	if err != nil {
		return "", err
	}
	user.Password = passwordHash

	if err := s.repo.Authorization.CreateUser(user); err != nil {
		return "", err
	}

	sendMailData := ActivationMailData{
		To: []string{user.Email},
		Link: activationLink,
	}
	sendMailDataJSON, err := json.Marshal(sendMailData)
	if err != nil {
		return "", err
	}
	if err := s.rabbitMQ.Publish(ACTIVATION, sendMailDataJSON); err != nil {
		return "", err
	}

	jwt, err := auth.GenerateToken(user.ID.String())
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func (s *AuthService) SignIn(user *model.User) (string, error) {
	userDB, err := s.repo.User.FindByUsernameOrEmail(user.Username, user.Email)
	if err != nil {
		return "", err
	}
	if userDB == nil {
		return "", errInvalidCredentials
	}

	if ok := auth.VerifyPassword([]byte(userDB.Password), []byte(user.Password)); !ok {
		return "", errInvalidCredentials
	}

	token, err := auth.GenerateToken(userDB.ID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Activate(ctx context.Context, activationLink string) error {
	user, err := s.repo.User.FindByActivationLink(activationLink)
	if err != nil {
		return err
	}
	if user == nil {
		return errUserNotFound
	}

	if err := s.rdb.Del(ctx, userPrefix + user.ID.String()).Err(); err != nil {
		return err
	}

	fields := map[string]interface{}{
		"activated": true,
		"activation_link": nil,
	}
	return s.repo.User.UpdateByID(user.ID, fields)
}

func (s *AuthService) SaveResetPasswordToken(token *model.ResetPasswordToken) error {
	hash, err := hasher.NewHash(12)
	if err != nil {
		return err
	}
	token.Token = hash

	if err := s.repo.Authorization.SaveResetPasswordToken(token); err != nil {
		return err
	}

	sendMailData := ResetPasswordTokenMailData{
		To: []string{token.UserEmail},
		Token: token.Token,
	}
	sendMailDataJSON, err := json.Marshal(sendMailData)
	if err != nil {
		return err
	}
	return s.rabbitMQ.Publish(RESET_TOKENS, sendMailDataJSON)
}

func (s *AuthService) FindResetPasswordToken(token string) (*model.ResetPasswordToken, error) {
	tokenDB, err := s.repo.Authorization.FindResetPasswordToken(token)
	if err != nil {
		return nil, err
	}

	return tokenDB, nil
}

func (s *AuthService) ResetPassword(token string, newPassword string) error {
	user, err := s.repo.User.FindByToken(token)
	if err != nil {
		return err
	}

	tokenDB, err := s.repo.Authorization.FindResetPasswordToken(token)
	if err != nil {
		return err
	}

	if expired := exp.IsExpired(tokenDB.Expiry); expired {
		if err := s.repo.Authorization.DeleteResetPasswordToken(token); err != nil {
			return err
		}
		return errTokenHasExpired
	}
	
	newHash, err := hasher.HashPassword([]byte(newPassword))
	if err != nil {
		return err
	}

	updateFields := map[string]interface{}{
		"password": newHash,
	}
	if err := s.repo.User.UpdateByID(user.ID, updateFields); err != nil {
		return err
	}

	return s.repo.Authorization.DeleteResetPasswordToken(token)
}
