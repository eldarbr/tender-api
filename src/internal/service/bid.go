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
	ErrBidCanceled     = errors.New("the bid is canceled and can't be changed")
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
	if b.AuthorType == model.AuthorTypeOrganization {
		idIsPresent, err := s.organizationRepo.GetOrganizationPresent(b.AuthorID)
		if err != nil {
			return err
		}
		if !idIsPresent {
			return ErrNoOrganization
		}
	} else if b.AuthorType == model.AuthorTypeUser {
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
	if currentBid.Status == model.BidPublished {
		return currentBid.Status, nil
	}
	err = authorizeUserForBid(username, currentBid, s.employeeRepo, s.organizationResponsibleRepo)
	if err != nil {
		return "", err
	}
	return currentBid.Status, nil
}

func (s *BidService) UpdateBidStatus(b *model.Bid, username string) error {
	currentBid, err := s.bidRepo.GetLastBidByID(b.ID)
	if err != nil {
		return err
	}
	err = authorizeUserForBid(username, currentBid, s.employeeRepo, s.organizationResponsibleRepo)
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
	err = authorizeUserForBid(username, currentBid, s.employeeRepo, s.organizationResponsibleRepo)
	if err != nil {
		return nil, err
	}
	if currentBid.Status == model.BidCanceled {
		return nil, ErrBidCanceled
	}
	return s.bidRepo.PatchBid(currentBid.ID, update)
}

func (s *BidService) RollbackBid(bidID uuid.UUID, username string, version int) (*model.Bid, error) {
	currentBid, err := s.bidRepo.GetLastBidByID(bidID)
	if err != nil {
		return nil, err
	}
	err = authorizeUserForBid(username, currentBid, s.employeeRepo, s.organizationResponsibleRepo)
	if err != nil {
		return nil, err
	}
	if currentBid.Status == model.BidCanceled {
		return nil, ErrBidCanceled
	}
	return s.bidRepo.RollbackBid(bidID, version)
}

func (s *BidService) LeaveFeedback(username string, bidID uuid.UUID, feedback string) (*model.Bid, error) {
	userID, err := s.employeeRepo.GetEmployeeIDByUsername(username)
	if err != nil {
		return nil, err
	}
	err = authorizeTenderResponsibleForBid(*userID, bidID, s.tenderRepo, s.bidRepo, s.organizationResponsibleRepo)
	if err != nil {
		return nil, err
	}
	return s.bidRepo.LeaveReview(bidID, feedback)
}

func authorizeUserForBid(username string, bid *model.Bid,
	employeeRepo *repository.EmployeeRepository,
	organizationResponsibleRepo *repository.OrganizationResponsibleRepository) error {
	employeeID, err := employeeRepo.GetEmployeeIDByUsername(username)
	if err != nil {
		return err
	}
	if bid.AuthorType == model.AuthorTypeUser && bid.AuthorID == *employeeID {
		return nil
	} else if bid.AuthorType == model.AuthorTypeOrganization {
		isResponsible, err := organizationResponsibleRepo.GetIfEmployeeIsResponsible(employeeID, &bid.AuthorID)
		if err != nil {
			return err
		}
		if isResponsible {
			return nil
		}
	}
	return ErrNotResponsible
}

func authorizeTenderResponsibleForBid(userID, bidID uuid.UUID,
	tenderRepo *repository.TenderRepository, bidRepo *repository.BidRepository,
	organizationResponsibleRepo *repository.OrganizationResponsibleRepository) error {
	currenctBid, err := bidRepo.GetLastBidByID(bidID)
	if err != nil {
		return err
	}
	if currenctBid.Status != model.BidPublished {
		return ErrNoBid
	}
	// authorize tender responsible
	tender, err := tenderRepo.GetLastTenderByID(currenctBid.TenderID)
	if err != nil {
		return err
	}
	isResponsible, err := organizationResponsibleRepo.GetIfEmployeeIsResponsible(&userID, &tender.OrganizationID)
	if err != nil {
		return err
	}
	if !isResponsible {
		return ErrNotResponsible
	}
	return nil
}

func (s *BidService) GetTenderReviewsOnUser(tenderID uuid.UUID, authorUsername, requesterUsername string,
	limit, offset int) ([]model.BidReview, error) {

	requesterID, err := s.employeeRepo.GetEmployeeIDByUsername(requesterUsername)
	if err != nil {
		return nil, err
	}
	tender, err := s.tenderRepo.GetLastTenderByID(tenderID)
	if err != nil {
		return nil, err
	}
	isResponsible, err := s.organizationResponsibleRepo.GetIfEmployeeIsResponsible(requesterID, &tender.OrganizationID)
	if err != nil {
		return nil, err
	}
	if !isResponsible {
		return nil, ErrNotResponsible
	}
	bidUserID, err := s.employeeRepo.GetEmployeeIDByUsername(authorUsername)
	if err != nil {
		return nil, err
	}

	return s.bidRepo.GetTenderReviewsOnUser(tenderID, *bidUserID, limit, offset)
}
