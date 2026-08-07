package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type LogService struct {
	repo repository.Log
}

func NewLogService(repo repository.Log) *LogService {
	return &LogService{repo: repo}
}

//func (s *LogService) Save(infoRu, inforTm, logType string) (int, error) {
//	return s.repo.Save(infoRu, inforTm, logType)
//}

func (s *LogService) GetAll(language string, pageSize, pageNumber int, logType string) (*models.LogPage, error) {
	return s.repo.GetAll(language, pageSize, pageNumber, logType)
}

func (s *LogService) GetLogTypes(language string) (*[]string, error) {
	return s.repo.GetLogTypes(language)
}
