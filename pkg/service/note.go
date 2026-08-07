package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type NoteService struct {
	repo repository.Note
}

func NewNoteService(repo repository.Note) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) Create(note models.CreateNote) error {
	return s.repo.Create(note)
}

func (s *NoteService) GetAll(pageSize, pageNumber int) (*models.NotePage, error) {
	return s.repo.GetAll(pageSize, pageNumber)
}

func (s *NoteService) Update(Id int, note models.CreateNote) error {
	return s.repo.Update(Id, note)
}

func (s *NoteService) Delete(Id int) error {
	return s.repo.Delete(Id)
}

func (s *NoteService) CreatePointNote(note models.CreatePointNote) error {
	return s.repo.CreatePointNote(note)
}

func (s *NoteService) GetAllPointNote(pageSize, pageNumber, tradePointId int, code, language string) (*models.PointNotePage, error) {
	return s.repo.GetAllPointNote(pageSize, pageNumber, tradePointId, language, code)
}

func (s *NoteService) UpdatePointNote(Id int, pointNote models.CreatePointNote) error {
	return s.repo.UpdatePointNote(Id, pointNote)
}

func (s *NoteService) DeletePointNote(Id int) error {
	return s.repo.DeletePointNote(Id)
}
