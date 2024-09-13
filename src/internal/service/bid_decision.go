package service

import (
	"avito-back-test/internal/model"
	"avito-back-test/internal/repository"
	"database/sql"
	"sync"

	"github.com/google/uuid"
)

var bidDecisionMutex sync.Mutex

type BidDecisionService struct {
	bidDecisionRepo         *repository.BidDecisionRepository
	bidRepo                 *repository.BidRepository
	tenderRepo              *repository.TenderRepository
	employeeRepo            *repository.EmployeeRepository
	organizationResponsRepo *repository.OrganizationResponsibleRepository
}

func NewBidDecisionService() *BidDecisionService {
	bidDesRepo := repository.NewBidDecisionRepository()
	bidRepo := repository.NewBidRepository()
	tenderRepo := repository.NewTenderRepository()
	emploRepo := repository.NewEmployeeRepository()
	orgRespRepo := repository.NewOrganizationResponsibleRepository()
	return &BidDecisionService{
		bidDecisionRepo:         bidDesRepo,
		bidRepo:                 bidRepo,
		tenderRepo:              tenderRepo,
		employeeRepo:            emploRepo,
		organizationResponsRepo: orgRespRepo,
	}
}

func (s *BidDecisionService) SubmitDecision(bidID uuid.UUID, username string, decision string) (*model.Bid, error) {
	bidDecisionMutex.Lock()
	defer bidDecisionMutex.Unlock()

	userID, err := s.employeeRepo.GetEmployeeIDByUsername(username)
	if err != nil {
		return nil, err
	}
	err = authorizeTenderResponsibleForBid(*userID, bidID, s.tenderRepo, s.bidRepo, s.organizationResponsRepo)
	if err != nil {
		return nil, err
	}
	currentBid, err := s.bidRepo.GetLastBidByID(bidID)
	if err != nil {
		return nil, err
	}

	err = s.bidDecisionRepo.WithTransaction(func(tx *sql.Tx) error {

		err := s.bidDecisionRepo.TxInsertUpdateDecision(tx, bidID, *userID, decision)
		if err != nil {
			return err
		}

		pro, contra, err := s.bidDecisionRepo.TxCountDecisions(tx, bidID)
		if err != nil {
			return err
		}
		if contra > 0 {
			err = s.bidRepo.TxSetBidStatus(tx, bidID, model.BidCanceled)
			if err != nil {
				return err
			}
			return nil
		}
		organizationRespCount, err := s.organizationResponsRepo.TxGetResponsibleCountByEmployee(tx, *userID)
		if err != nil {
			return err
		}
		quorum := min(organizationRespCount, 3)
		if pro > quorum {
			err = s.tenderRepo.TxUpdateTenderStatus(tx, currentBid.TenderID, model.TenderClosed)
		}
		return err
	})

	if err != nil {
		return nil, err
	}
	return s.bidRepo.GetLastBidByID(bidID)
}
