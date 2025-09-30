package service

import (
	"errors"
	"time"

	"github.com/fazel/notebooq/internal/models"
	"github.com/fazel/notebooq/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *repository.UserRepo
	jwtSecret string
	accessExp time.Duration
}

func NewUserService(r *repository.UserRepo, jwtSecret string, exp time.Duration) *UserService {
	return &UserService{repo: r, jwtSecret: jwtSecret, accessExp: exp}
}

func (s *UserService) Signup(username, password string) (*models.User, error) {
	existing, _ := s.repo.FindByUsername(username)
	if existing != nil {
		return nil, errors.New("username taken")
	}
	h, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u := &models.User{Username: username, Password: string(h)}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}
func (s *UserService) GetByID(id uint) (*models.User, error) {
	return s.repo.GetByID(id)
}
func (s *UserService) Login(username, password string) (string, *models.User, error) {
	u, err := s.repo.FindByUsername(username)
	if err != nil || u == nil {
		return "", nil, errors.New("invalid credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// create token
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": u.ID,
		"exp": time.Now().Add(s.accessExp).Unix(),
	})
	signed, err := t.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", nil, err
	}
	return signed, u, nil
}
