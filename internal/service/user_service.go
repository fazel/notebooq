package service

import (
	"errors"
	"math/rand"
	"regexp"
	"time"

	"github.com/fazel/notebooq/internal/models"
	"github.com/fazel/notebooq/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo        *repository.UserRepo
	JWTSecret   string
	TokenExpire time.Duration
}

func NewUserService(r *repository.UserRepo, secret string, expire time.Duration) *UserService {
	return &UserService{
		repo:        r,
		JWTSecret:   secret,
		TokenExpire: expire,
	}
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *UserService) GenerateJWT(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.TokenExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JWTSecret))
}
func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return errors.New("password must contain at least one number")
	}
	if !regexp.MustCompile(`[!@#~$%^&*()+|_]`).MatchString(password) {
		return errors.New("password must contain at least one special character")
	}
	return nil
}

func (s *UserService) CreateUser(username, password, email, code string) (*models.User, error) {
	if err := validatePassword(password); err != nil {
		return nil, err
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u := &models.User{
		Username:   username,
		Password:   string(hashed),
		Email:      email,
		IsVerified: false,
		VerifyCode: code,
	}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func GenerateCode() string {
	const letters = "0123456789"
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, 6)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}

func (s *UserService) Update(u *models.User) error {
	return s.repo.Update(u)
}

func (s *UserService) GetByUsername(username string) (*models.User, error) {
	return s.repo.GetByUsername(username)
}

func (s *UserService) GetByID(id uint) (*models.User, error) {
	return s.repo.GetByID(id)
}
func (s *UserService) Signup(username, email, password string) (*models.User, string, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	code := generateCode(6)
	u := &models.User{
		Username:   username,
		Email:      email,
		Password:   string(hash),
		IsVerified: false,
		VerifyCode: code,
	}
	s.repo.Create(u)
	return u, code, nil
}

// Generate random numeric code
func generateCode(length int) string {
	const digits = "0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = digits[int(b[i])%len(digits)]
	}
	return string(b)
}
