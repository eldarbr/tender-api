package service

import (
	"avito-back-test/internal/model"
	"avito-back-test/internal/repository"
)

type TenderService struct {
	repo *repository.TenderRepository
}

func NewTenderService() *TenderService {
	repo := repository.NewTenderRepository()
	return &TenderService{
		repo: repo,
	}
}

func (s *TenderService) GetTenders() ([]model.Tender, error) {
	return s.repo.GetAllTenders()
}
