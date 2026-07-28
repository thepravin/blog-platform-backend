package services

import (
	"blog_platform/internal/auth"
	"blog_platform/internal/models"
	"blog_platform/internal/repositories"
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(r *repositories.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) Register(username, email, password string) (*models.User, error) {
	u := &models.User{
		ID:          uuid.New(),
		UserName:    username,
		Email:       email,
		Password:    password,
		DisplayName: username,
		Role:        "reader",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	hashed, err := auth.HashedPassword(password)
	if err != nil {
		return nil, err
	}
	u.Password = hashed

	if err := s.repo.Create(u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UserService) Login(email, password string) (*models.User, string, error) {
	u, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, "", errors.New("Invalid Email or Password")
	}

	if !auth.CheckPassword(u.Password, password) {
		return nil, "", errors.New("Invalid Email or Password")
	}

	token, _ := auth.GenerateJWT(u.ID.String(), u.Email, u.Role)
	return u, token, nil
}

func (s *UserService) GetProfile(id string) (*models.User, error) {
	u, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return u, nil
}
