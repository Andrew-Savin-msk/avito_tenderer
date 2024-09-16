package bidservice

import (
	"errors"

	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/domain/models"
	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/services"
	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/store"
	"github.com/sirupsen/logrus"
)

type Bider struct {
	ts     store.Tenders
	bs     store.Bids
	rs     store.Responsibles
	logger *logrus.Entry
}

func New(tenderStore store.Tenders, bidStorage store.Bids, responsiblesStore store.Responsibles, log *logrus.Logger) *Bider {
	logger := log.WithFields(logrus.Fields{
		"service": "bider",
	})

	return &Bider{
		ts:     tenderStore,
		bs:     bidStorage,
		rs:     responsiblesStore,
		logger: logger,
	}
}

func (b *Bider) Create(bid *models.Bid) (*models.Bid, error) {
	var orgId string
	var err error
	if bid.AuthorType == "Organization" {
		orgId = bid.AuthorId
	} else {
		orgId, err = b.rs.ResponcibleForOrg(bid.AuthorId)
		if err != nil {
			if errors.Is(err, store.ErrConnClosed) {
				return nil, services.ErrServiceDatabaseDisconnected
			}
			if errors.Is(err, store.ErrRecordNotFound) {
				return nil, services.ErrNoSuchUser
			}
			b.logger.Errorf("unexpected error: %s on method Create", err)
			return nil, err
		}
	}

	tenderCondition, err := b.ts.GetCondition(bid.TenderId, store.Latest)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		if errors.Is(err, store.ErrRecordNotFound) {
			return nil, services.ErrNoSuchTender
		}
		b.logger.Errorf("unexpected error: %s on method GetCondition", err)
		return nil, err
	}

	if tenderCondition.OrgId != "Published" && tenderCondition.OrgId != orgId {
		return nil, services.ErrNoPermitions
	}

	data, err := b.bs.Create(bid, orgId)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		b.logger.Errorf("unexpected error: %s on method Create", err)
		return nil, err
	}
	return data, nil
}

func (b *Bider) GetByName(limit, offset int64, username string) ([]*models.Bid, error) {
	// получить user_id
	userId, err := b.rs.GetUserId(username)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		if errors.Is(err, store.ErrUserNotFound) {
			return nil, services.ErrNoSuchUser
		}
		b.logger.Errorf("unexpected error: %s on method GetUserId", err)
		return nil, err
	}

	result, err := b.bs.GetUserList(limit, offset, userId)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		b.logger.Errorf("unexpected error: %s on method GetLimitedList", err)
		return nil, err
	}

	return result, nil
}

func (b *Bider) GetTenderBids(limit, offset int64, tenderId, username string) ([]*models.Bid, error) {
	userOrgId, err := b.rs.GetOrgId(username)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		if errors.Is(err, store.ErrUserNotFound) {
			return nil, services.ErrNoSuchUser
		}
		if errors.Is(err, store.ErrRecordNotFound) {
			return nil, services.ErrNoPermitions
		}
		b.logger.Errorf("unexpected error: %s on method GetOrgId", err)
		return nil, err
	}

	tenderCondition, err := b.ts.GetCondition(tenderId, store.Latest)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		if errors.Is(err, store.ErrRecordNotFound) {
			return nil, services.ErrNoSuchTender
		}
		b.logger.Errorf("unexpected error: %s on method GetCondition", err)
		return nil, err
	}

	if tenderCondition.Status != "Published" && tenderCondition.OrgId != userOrgId {
		return nil, services.ErrNoPermitions
	}

	result, err := b.bs.GetTenderList(limit, offset, tenderId, userOrgId)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		b.logger.Errorf("unexpected error: %s on method GetTenderList", err)
		return nil, err
	}
	return result, nil
}

func (b *Bider) GetStat(bidId, username string) (string, error) {
	userOrgId, err := b.rs.GetOrgId(username)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return "", services.ErrServiceDatabaseDisconnected
		}
		if errors.Is(err, store.ErrUserNotFound) {
			return "", services.ErrNoSuchUser
		}
		if errors.Is(err, store.ErrRecordNotFound) {
			return "", services.ErrNoPermitions
		}
		b.logger.Errorf("unexpected error: %s on method GetUserId", err)
		return "", err
	}

	bidCondition, err := b.bs.GetCondition(bidId, store.Latest)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return "", services.ErrServiceDatabaseDisconnected
		}
		if errors.Is(err, store.ErrRecordNotFound) {
			return "", services.ErrNoSuchBid
		}
		b.logger.Errorf("unexpected error: %s on method GetCondition", err)
		return "", err
	}

	if bidCondition.Status != "Published" {
		if bidCondition.AuthorType == "Organization" {
			if bidCondition.AuthorId != userOrgId {
				return "", services.ErrNoPermitions
			}
			return bidCondition.Status, nil
		} else {
			bidOrgId, err := b.rs.ResponcibleForOrg(bidCondition.AuthorId)
			if err != nil {
				if errors.Is(err, store.ErrConnClosed) {
					return "", services.ErrServiceDatabaseDisconnected
				}
				if errors.Is(err, store.ErrRecordNotFound) {
					return "", services.ErrNoSuchBid
				}
				b.logger.Errorf("unexpected error: %s on method ResponcibleForOrg", err)
				return "", err
			}

			if bidOrgId != userOrgId {
				return "", services.ErrNoPermitions
			}
		}
	} else {
		bidOrgId, err := b.rs.ResponcibleForOrg(bidCondition.AuthorId)
		if err != nil {
			if errors.Is(err, store.ErrConnClosed) {
				return "", services.ErrServiceDatabaseDisconnected
			}
			if errors.Is(err, store.ErrRecordNotFound) {
				return "", services.ErrNoSuchBid
			}
			b.logger.Errorf("unexpected error: %s on method ResponcibleForOrg", err)
			return "", err
		}

		tenderCond, err := b.ts.GetCondition(bidCondition.TenderId, store.Latest)
		if err != nil {
			if errors.Is(err, store.ErrConnClosed) {
				return "", services.ErrServiceDatabaseDisconnected
			}
			if errors.Is(err, store.ErrRecordNotFound) {
				return "", services.ErrNoSuchTender
			}
			b.logger.Errorf("unexpected error: %s on method ResponcibleForOrg", err)
			return "", err
		}

		if userOrgId != bidOrgId && userOrgId != tenderCond.OrgId {
			return "", services.ErrNoPermitions
		}
	}

	return bidCondition.Status, nil
}
func (b *Bider) ChangeStat(bidId, status, username string) (*models.Bid, error) {
	bidCondition, err := b.bs.GetCondition(bidId, store.Latest)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		if errors.Is(err, store.ErrRecordNotFound) {
			return nil, services.ErrNoSuchBid
		}
		b.logger.Errorf("unexpected error: %s on method GetCondition", err)
		return nil, err
	}

	userOrgId, err := b.rs.GetOrgId(username)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		if errors.Is(err, store.ErrUserNotFound) {
			return nil, services.ErrNoSuchUser
		}
		if errors.Is(err, store.ErrRecordNotFound) {
			return nil, services.ErrNoPermitions
		}
		b.logger.Errorf("unexpected error: %s on method GetOrgId", err)
		return nil, err
	}

	if bidCondition.AuthorType == "Organization" {
		if bidCondition.AuthorId != userOrgId {
			return nil, services.ErrNoPermitions
		}
	} else {
		bidOrgId, err := b.rs.ResponcibleForOrg(bidCondition.AuthorId)
		if err != nil {
			if errors.Is(err, store.ErrConnClosed) {
				return nil, services.ErrServiceDatabaseDisconnected
			}
			if errors.Is(err, store.ErrRecordNotFound) {
				return nil, services.ErrNoPermitions
			}
			b.logger.Errorf("unexpected error: %s on method ResponcibleForOrg", err)
			return nil, err
		}

		if bidOrgId != userOrgId {
			return nil, services.ErrNoPermitions
		}
	}

	if bidCondition.Status == status {
		return bidCondition, nil
	}

	bidCondition.Status = status
	bidCondition.Version += 1

	result, err := b.bs.UpdateCondition(bidCondition)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		b.logger.Errorf("unexpected error: %s on method UpdateCondition", err)
		return nil, err
	}
	return result, nil
}

func (b *Bider) Edit(bid *models.Bid, bidId, username string) (*models.Bid, error) {
	bidCondition, err := b.bs.GetCondition(bidId, store.Latest)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		if errors.Is(err, store.ErrRecordNotFound) {
			return nil, services.ErrNoSuchBid
		}
		b.logger.Errorf("unexpected error: %s on method GetCondition", err)
		return nil, err
	}

	userOrgId, err := b.rs.GetOrgId(username)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		if errors.Is(err, store.ErrUserNotFound) {
			return nil, services.ErrNoSuchUser
		}
		if errors.Is(err, store.ErrRecordNotFound) {
			return nil, services.ErrNoPermitions
		}
		b.logger.Errorf("unexpected error: %s on method GetOrgId", err)
		return nil, err
	}

	if bidCondition.AuthorType == "Organization" {
		if bidCondition.AuthorId != userOrgId {
			return nil, services.ErrNoPermitions
		}
	} else {
		bidOrgId, err := b.rs.ResponcibleForOrg(bidCondition.AuthorId)
		if err != nil {
			if errors.Is(err, store.ErrConnClosed) {
				return nil, services.ErrServiceDatabaseDisconnected
			}
			if errors.Is(err, store.ErrRecordNotFound) {
				return nil, services.ErrNoPermitions
			}
			b.logger.Errorf("unexpected error: %s on method ResponcibleForOrg", err)
			return nil, err
		}

		if bidOrgId != userOrgId {
			return nil, services.ErrNoPermitions
		}
	}

	var count int
	if bid.Name != "" && bidCondition.Name != bid.Name {
		bidCondition.Name = bid.Name
		count++
	}
	if bid.Description != "" && bidCondition.Description != bid.Description {
		bidCondition.Description = bid.Description
		count++
	}
	if count == 0 {
		return bidCondition, nil
	}

	bidCondition.Version += 1

	result, err := b.bs.UpdateCondition(bidCondition)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		b.logger.Errorf("unexpected error: %s on method UpdateCondition", err)
		return nil, err
	}
	return result, nil
}
func (b *Bider) Sumbit(bidId, decision, username string) (*models.Bid, error) {
	panic("unimplemented")
}

func (b *Bider) Rollback(bidId string, version int64, username string) (*models.Bid, error) {
	bidCondition, err := b.bs.GetCondition(bidId, version)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		if errors.Is(err, store.ErrRecordNotFound) {
			return nil, services.ErrNoSuchBid
		}
		b.logger.Errorf("unexpected error: %s on method GetCondition", err)
		return nil, err
	}

	userOrgId, err := b.rs.GetOrgId(username)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		if errors.Is(err, store.ErrUserNotFound) {
			return nil, services.ErrNoSuchUser
		}
		if errors.Is(err, store.ErrRecordNotFound) {
			return nil, services.ErrNoPermitions
		}
		b.logger.Errorf("unexpected error: %s on method GetOrgId", err)
		return nil, err
	}

	if bidCondition.AuthorType == "Organization" {
		if bidCondition.AuthorId != userOrgId {
			return nil, services.ErrNoPermitions
		}
	} else {
		bidOrgId, err := b.rs.ResponcibleForOrg(bidCondition.AuthorId)
		if err != nil {
			if errors.Is(err, store.ErrConnClosed) {
				return nil, services.ErrServiceDatabaseDisconnected
			}
			if errors.Is(err, store.ErrRecordNotFound) {
				return nil, services.ErrNoPermitions
			}
			b.logger.Errorf("unexpected error: %s on method ResponcibleForOrg", err)
			return nil, err
		}

		if bidOrgId != userOrgId {
			return nil, services.ErrNoPermitions
		}
	}

	latestVersion, err := b.bs.GetBidLatestVersion(bidId)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		if errors.Is(err, store.ErrRecordNotFound) {
			return nil, services.ErrNoSuchBid
		}
		b.logger.Errorf("unexpected error: %s on method GetBidLatestVersion", err)
		return nil, err
	}

	bidCondition.Version = latestVersion + 1

	result, err := b.bs.UpdateCondition(bidCondition)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		b.logger.Errorf("unexpected error: %s on method UpdateCondition", err)
		return nil, err
	}
	return result, nil
}
