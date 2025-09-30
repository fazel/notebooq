package repository

import (
	"github.com/fazel/notebooq/internal/models"
	"gorm.io/gorm"
)

type NoteRepo struct {
	db *gorm.DB
}

func NewNoteRepo(db *gorm.DB) *NoteRepo { return &NoteRepo{db: db} }

func (r *NoteRepo) Create(n *models.Note) error {
	return r.db.Create(n).Error
}

func (r *NoteRepo) ListByUser(userID uint) ([]models.Note, error) {
	var notes []models.Note
	if err := r.db.Where("user_id = ?", userID).Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *NoteRepo) GetByID(id uint) (*models.Note, error) {
	var n models.Note
	if err := r.db.First(&n, id).Error; err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *NoteRepo) Update(n *models.Note) error {
	return r.db.Save(n).Error
}

func (r *NoteRepo) Delete(id uint) error {
	return r.db.Delete(&models.Note{}, id).Error
}
