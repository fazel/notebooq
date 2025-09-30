package service

import (
	"github.com/fazel/notebooq/internal/models"
	"github.com/fazel/notebooq/internal/repository"
)

type NoteService struct {
	repo *repository.NoteRepo
}

func NewNoteService(r *repository.NoteRepo) *NoteService { return &NoteService{repo: r} }

func (s *NoteService) Create(n *models.Note) error { return s.repo.Create(n) }
func (s *NoteService) ListByUser(userID uint) ([]models.Note, error) {
	return s.repo.ListByUser(userID)
}
func (s *NoteService) GetByID(id uint) (*models.Note, error) { return s.repo.GetByID(id) }
func (s *NoteService) Update(n *models.Note) error           { return s.repo.Update(n) }
func (s *NoteService) Delete(id uint) error                  { return s.repo.Delete(id) }
