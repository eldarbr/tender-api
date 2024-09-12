package service

import (
	"avito-back-test/internal/model"
	"avito-back-test/internal/repository"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrWrongAuthorType = errors.New("author type not supported")
	ErrNoOrganization  = repository.ErrNoOrganization
	ErrNoBid           = repository.ErrNoBid
)

type BidService struct {
	bidRepo                     *repository.BidRepository
	tenderRepo                  *repository.TenderRepository
	employeeRepo                *repository.EmployeeRepository
	organizationResponsibleRepo *repository.OrganizationResponsibleRepository
	organizationRepo            *repository.OrganizationRepository
}

func NewBidService() *BidService {
	bidRepo := repository.NewBidRepository()
	tenderRepo := repository.NewTenderRepository()
	organizationResponsibleRepo := repository.NewOrganizationResponsibleRepository()
	employeeRepo := repository.NewEmployeeRepository()
	orgagizationRepo := repository.NewOrganizationRepository()
	return &BidService{
		bidRepo:                     bidRepo,
		tenderRepo:                  tenderRepo,
		organizationResponsibleRepo: organizationResponsibleRepo,
		employeeRepo:                employeeRepo,
		organizationRepo:            orgagizationRepo,
	}
}

func (s *BidService) InsertNewBid(b *model.Bid) error {
	if b.AuthorType == "Organization" {
		idIsPresent, err := s.organizationRepo.GetOrganizationPresent(b.AuthorID)
		if err != nil {
			return err
		}
		if !idIsPresent {
			return ErrNoOrganization
		}
	} else if b.AuthorType == "User" {
		idIsPresent, err := s.employeeRepo.GetEmployeePresent(b.AuthorID)
		if err != nil {
			return err
		}
		if !idIsPresent {
			return ErrNoEmployee
		}
		// check if the user is responsible
		if _, err := s.employeeRepo.GetEmployeeRespOrganization(b.AuthorID); err != nil {
			return err
		}
	} else {
		return ErrWrongAuthorType
	}
	_, err := s.tenderRepo.GetLastTenderByID(b.TenderID)
	if err != nil {
		return err
	}

	return s.bidRepo.InsertNewBid(b)
}

func (s *BidService) GetUserBids(username string, limit, offset int) ([]model.Bid, error) {
	// check username validity
	employeeID, err := s.employeeRepo.GetEmployeeIDByUsername(username)
	if err != nil {
		return nil, err
	}
	return s.bidRepo.GetUserBids(*employeeID, limit, offset)
}

func (s *BidService) GetBidsByTender(tenderID uuid.UUID, username string, limit, offset int) ([]model.Bid, error) {
	// check username validity
	employeeID, err := s.employeeRepo.GetEmployeeIDByUsername(username)
	if err != nil {
		return nil, err
	}
	// check if the user is responsible for the tender
	tender, err := s.tenderRepo.GetLastTenderByID(tenderID)
	if err != nil {
		return nil, err
	}
	isResponsible, err := s.organizationResponsibleRepo.GetIfEmployeeIsResponsible(employeeID, &tender.OrganizationID)
	if err != nil {
		return nil, err
	}
	if !isResponsible {
		return nil, ErrNotResponsible
	}
	return s.bidRepo.GetUserBids(*employeeID, limit, offset)
}

func (s *BidService) GetBidStatus(bidID uuid.UUID, username string) (string, error) {
	currentBid, err := s.bidRepo.GetLastBidByID(bidID)
	if err != nil {
		return "", err
	}
	// return immediately if the bid is public
	if currentBid.Status == "Published" {
		return currentBid.Status, nil
	}
	err = s.authorizeUserForBid(username, currentBid)
	if err != nil {
		return "", err
	}
	return currentBid.Status, nil
}

func (s *BidService) UpdateTenderStatus(b *model.Bid, username string) error {
	currentBid, err := s.bidRepo.GetLastBidByID(b.ID)
	if err != nil {
		return err
	}
	err = s.authorizeUserForBid(username, currentBid)
	if err != nil {
		return err
	}
	return s.bidRepo.UpdateBidStatus(b)
}

func (s *BidService) PatchBid(bidID uuid.UUID, username string, update *model.BidUpdate) (*model.Bid, error) {
	currentBid, err := s.bidRepo.GetLastBidByID(bidID)
	if err != nil {
		return nil, err
	}
	err = s.authorizeUserForBid(username, currentBid)
	if err != nil {
		return nil, err
	}
	return s.bidRepo.PatchBid(currentBid.ID, update)
}

func (s *BidService) authorizeUserForBid(username string, bid *model.Bid) error {
	employeeID, err := s.employeeRepo.GetEmployeeIDByUsername(username)
	if err != nil {
		return err
	}
	if bid.AuthorType == "User" && bid.AuthorID == *employeeID {
		return nil
	} else if bid.AuthorType == "Organization" {
		isResponsible, err := s.organizationResponsibleRepo.GetIfEmployeeIsResponsible(employeeID, &bid.AuthorID)
		if err != nil {
			return err
		}
		if isResponsible {
			return nil
		}
	}
	return ErrNotResponsible
}
