package service

import (
	"avito-back-test/internal/model"
	"avito-back-test/internal/repository"
	"errors"
)

var ErrNotResponsible = errors.New("the employee is not responsible")

type TenderService struct {
	tenderRepo                  *repository.TenderRepository
	organizationResponsibleRepo *repository.OrganizationResponsibleRepository
	employeeRepo                *repository.EmployeeRepository
}

func NewTenderService() *TenderService {
	tenderRepo := repository.NewTenderRepository()
	organizationResponsibleRepo := repository.NewOrganizationResponsibleRepository()
	employeeRepo := repository.NewEmployeeRepository()
	return &TenderService{
		tenderRepo:                  tenderRepo,
		organizationResponsibleRepo: organizationResponsibleRepo,
		employeeRepo:                employeeRepo,
	}
}

func (s *TenderService) GetTenders() ([]model.Tender, error) {
	return s.tenderRepo.GetAllTenders()
}

func (s *TenderService) InsertNewTender(t *model.Tender) error {
	// Get id by username
	employeeId, err := s.employeeRepo.GetEmployeeIDByUsername(t.CreatorUsername)
	if err != nil {
		return err
	}
	// Check if the employee is responsible
	isResponsible, err := s.organizationResponsibleRepo.GetIfEmployeeIsResponsible(employeeId, &t.OrganizationID)
	if err != nil {
		return err
	}
	if !isResponsible {
		return ErrNotResponsible
	}
	return s.tenderRepo.InsertNewTender(t)
}
