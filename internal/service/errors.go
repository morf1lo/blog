package service

import "errors"

var (
	errUserIsAlreadyExists = errors.New("user with this email is already exists")
	errInvalidCredentials = errors.New("invalid credentials")
	errUserNotFound = errors.New("user not found")
	errTokenHasExpired = errors.New("token has expired")
	errFileIsNotAnImage = errors.New("file is not an image")
	errNoAccess = errors.New("you have no access")
)
