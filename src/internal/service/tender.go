package service

import (
	"avito-back-test/internal/model"
	"avito-back-test/internal/repository"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNotResponsible = errors.New("the employee is not responsible")
	ErrNoEmployee     = repository.ErrNoEmployee
	ErrNoTender       = repository.ErrNoTender
)

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
	return s.tenderRepo.GetAllPublicTenders(limit, offset)
}

func (s *TenderService) GetTendersOfService(service string, limit, offset int) ([]model.Tender, error) {
	return s.tenderRepo.GetPublicTendersOfService(service, limit, offset)
}

func (s *TenderService) InsertNewTender(t *model.Tender, username string) error {
	// Get id by username
	employeeId, err := s.employeeRepo.GetEmployeeIDByUsername(username)
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
	employeeID, err := s.employeeRepo.GetEmployeeIDByUsername(username)
	if err != nil {
		return nil, err
	}
	return s.tenderRepo.GetUserTenders(*employeeID, limit, offset)
}

func (s *TenderService) GetTenderStatus(tenderID uuid.UUID, username *string) (string, error) {
	currentTender, err := s.tenderRepo.GetTenderByID(tenderID)
	if err != nil {
		return "", err
	}
	// return immediately if the tender is public
	if currentTender.Status == "Published" {
		return currentTender.Status, nil
	}
	// otherwise (not public) return responsibility error if no username is provided
	if username == nil {
		return "", ErrNotResponsible
	}
	// get user id by username
	employeeId, err := s.employeeRepo.GetEmployeeIDByUsername(*username)
	if err != nil {
		return "", err
	}
	// Check if the employee is responsible
	isResponsible, err := s.organizationResponsibleRepo.GetIfEmployeeIsResponsible(employeeId, &currentTender.OrganizationID)
	if err != nil {
		return "", err
	}
	if !isResponsible {
		return "", ErrNotResponsible
	}
	// at this point, the user is responsible and can see the response
	return currentTender.Status, nil
}

func (s *TenderService) UpdateTenderStatus(t *model.Tender, username string) error {
	// Get id by username
	employeeId, err := s.employeeRepo.GetEmployeeIDByUsername(username)
	if err != nil {
		return err
	}
	currentTender, err := s.tenderRepo.GetTenderByID(t.ID)
	if err != nil {
		return err
	}
	// Check if the employee is responsible
	isResponsible, err := s.organizationResponsibleRepo.GetIfEmployeeIsResponsible(employeeId, &currentTender.OrganizationID)
	if err != nil {
		return err
	}
	if !isResponsible {
		return ErrNotResponsible
	}
	return s.tenderRepo.UpdateTenderStatus(t)
}
