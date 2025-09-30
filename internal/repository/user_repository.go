package repository

import (
	"errors"

	"github.com/fazel/notebooq/internal/models"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(u *models.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepo) Update(u *models.User) error {
	return r.db.Save(u).Error
}

func (r *UserRepo) GetByUsername(username string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByID(id uint) (*models.User, error) {
	var u models.User
	if err := r.db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &u, nil
}
