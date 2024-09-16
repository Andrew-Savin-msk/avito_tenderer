package bidservice

import (
	"errors"

	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/domain/models"
	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/services"
	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/store"
)

func (b *Bider) AddFeedback(bidId, bidFeedback, username string) (*models.Bid, error) {
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

	userOrgId, err := b.rs.ResponcibleForOrg(userId)
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

	tenderOrgId, err := b.ts.GetOrgIdByBidId(bidId)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		if errors.Is(err, store.ErrUserNotFound) {
			return nil, services.ErrNoSuchUser
		}
		if errors.Is(err, store.ErrRecordNotFound) {
			return nil, services.ErrNoSucnResource
		}
		b.logger.Errorf("unexpected error: %s on method GetOrgIdByBidId", err)
		return nil, err
	}

	if tenderOrgId != userOrgId {
		return nil, services.ErrNoPermitions
	}

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


	if bidCondition.Status != "Published" {
		return nil, services.ErrNoPermitions
	}

	err = b.bs.AddFeedback(bidId, userId, bidFeedback)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		b.logger.Errorf("unexpected error: %s on method AddFeedback", err)
		return nil, err
	}

	return bidCondition, nil
}

func (b *Bider) GetReviews(tenderId, authorUsername, requesterUsername string, limit, offset int64) ([]*models.Feedback, error) {
	reqUserId, err := b.rs.GetUserId(requesterUsername)
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

	reqUserOrgId, err := b.rs.ResponcibleForOrg(reqUserId)
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

	if tenderCondition.OrgId != reqUserOrgId {
		return nil, services.ErrNoPermitions
	}

	result, err := b.bs.GetFeedbacks(tenderId, authorUsername, limit, offset)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		b.logger.Errorf("unexpected error: %s on method GetFeedbacks", err)
		return nil, err
	}
	return result, nil
}
