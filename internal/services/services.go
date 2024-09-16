package services

import "github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/domain/models"

type Tenders interface {
	List(limit, offset int64, serviceType []string) ([]*models.Tender, error)
	Create(tnd *models.Tender, responcible *models.Responsible) (*models.Tender, error)
	GetByName(limit, offset int64, username string) ([]*models.Tender, error)
	GetStat(tenderId, username string) (string, error)
	ChangeStat(tenderId, status, username string) (*models.Tender, error)
	Edit(tnd *models.Tender, tenderid, username string) (*models.Tender, error)
	Rollback(tenderId, username string, version int64) (*models.Tender, error)
}

type Bids interface {
	Create(bid *models.Bid) (*models.Bid, error)
	GetByName(limit, offset int64, username string) ([]*models.Bid, error)
	GetTenderBids(limit, offset int64, tenderId, username string) ([]*models.Bid, error)
	GetStat(bidId, username string) (string, error)
	ChangeStat(bidId, status, username string) (*models.Bid, error)
	Edit(bid *models.Bid, bidId, username string) (*models.Bid, error)
	Sumbit(bidId, decision, username string) (*models.Bid, error)
	AddFeedback(bidId, bidFeedback, username string) (*models.Bid, error)
	Rollback(bidId string, version int64, username string) (*models.Bid, error)
	GetReviews(tenderId, authorUsername, requesterUsername string, limit, offset int64) ([]*models.Feedback, error)
}
