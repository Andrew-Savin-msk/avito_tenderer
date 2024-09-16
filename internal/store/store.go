package store

import "github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/domain/models"

const Latest = -1

type Tenders interface {
	GetLimitedList(limit, offset int64, servType []string) ([]*models.Tender, error)
	Create(tnd *models.Tender, resp *models.Responsible) (*models.Tender, error)
	GetUserTenders(limit, offset int64, username string) ([]*models.Tender, error)
	GetStatus(tenderId string) (string, error)
	GetTenderLatestVersion(tenderId, username string) (int64, error)
	GetCondition(tenderId string, version int64) (*models.Tender, error)
	UpdateCondition(newCondition *models.Tender) (*models.Tender, error)
	IsResponcibleFor(tenderId string, respUUIDs []string) error
	GetOrgIdByBidId(bidId string) (string, error)
}

type Responsibles interface {
	GetResponsibleUUID(responsibles *models.Responsible) (string, error)
	GetRespUUIDs(username string) ([]string, error)
	IsResponcible(responsible *models.Responsible) error
	IsUserExists(username string) error
	GetOrgId(username string) (string, error)
	ResponcibleForOrg(userId string) (string, error)
	GetUserId(username string) (string, error)
}

type Bids interface {
	Create(bid *models.Bid, orgId string) (*models.Bid, error)
	GetUserList(limit, offset int64, userId string) ([]*models.Bid, error)
	GetTenderList(limit, offset int64, tenderId, orgId string) ([]*models.Bid, error)
	GetCondition(bidId string, version int64) (*models.Bid, error)
	GetBidLatestVersion(bidId string) (int64, error)
	UpdateCondition(newCondition *models.Bid) (*models.Bid, error)
	AddFeedback(bidId, userId, feedback string) error
	GetFeedbacks(tenderId, authorUserId string, limit, offset int64) ([]*models.Feedback, error)
}
