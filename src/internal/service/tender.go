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

func (s *TenderService) GetTenders(limit, offset int) ([]model.Tender, error) {
	return s.tenderRepo.GetAllTenders(limit, offset)
}

func (s *TenderService) GetTendersOfService(service string, limit, offset int) ([]model.Tender, error) {
	return s.tenderRepo.GetTendersOfService(service, limit, offset)
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

func (s *TenderService) GetUserTenders(username string, limit, offset int) ([]model.Tender, error) {
	// check username validity
	_, err := s.employeeRepo.GetEmployeeIDByUsername(username)
	if err != nil {
		return nil, err
	}
	return s.tenderRepo.GetUserTenders(username, limit, offset)
}
