package tenderservice

import (
	"errors"

	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/domain/models"
	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/services"
	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/store"
	"github.com/sirupsen/logrus"
)

/*
	type Tenders interface {
		List(limit, offset int64, serviceType []string) ([]*models.Tender, error)
		Create(tnd *models.Tender, orgId, username string) (*models.Tender, error)
		GetByName(limit, offset int64, username string) ([]*models.Tender, error)
		GetStat(tenderId, username string) (string, error)
		ChangeStat(tenderId, status, username string) (*models.Tender, error)
		Edit(tnd *models.Tender, tenderid, username string) (*models.Tender, error)
		Rollback(tenderId, username string, version int64) (*models.Tender, error)
	}
*/

type Tender struct {
	ts     store.Tenders
	rs     store.Responsibles
	logger *logrus.Entry
}

func New(tenderStorage store.Tenders, responsiblesStore store.Responsibles, log *logrus.Logger) *Tender {
	logger := log.WithFields(logrus.Fields{
		"service": "tender",
	})

	return &Tender{
		ts:     tenderStorage,
		rs:     responsiblesStore,
		logger: logger,
	}
}

func (t *Tender) List(limit, offset int64, serviceType []string) ([]*models.Tender, error) {
	tenders, err := t.ts.GetLimitedList(limit, offset, serviceType)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		t.logger.Errorf("unexpected error: %s on method GetLimitedList", err)
		return nil, err
	}
	return tenders, nil
}

func (t *Tender) Create(tnd *models.Tender, responsible *models.Responsible) (*models.Tender, error) {
	err := t.rs.IsResponcible(responsible)
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
		t.logger.Errorf("unexpected error: %s on method IsResponcible", err)
		return nil, err
	}

	result, err := t.ts.Create(tnd, responsible)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		t.logger.Errorf("unexpected error: %s on method Create", err)
		return nil, err
	}
	return result, nil
}

func (t *Tender) GetByName(limit, offset int64, username string) ([]*models.Tender, error) {
	err := t.rs.IsUserExists(username)
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
		t.logger.Errorf("unexpected error: %s on method IsUserExists", err)
		return nil, err
	}

	tenders, err := t.ts.GetUserTenders(limit, offset, username)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		t.logger.Errorf("unexpected error: %s on method GetUserTenders", err)
		return nil, err
	}
	return tenders, nil
}

func (t *Tender) GetStat(tenderId, username string) (string, error) {
	orgId, err := t.rs.GetOrgId(username)
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
		t.logger.Errorf("unexpected error: %s on method GetOrgId", err)
		return "", err
	}

	tenderCondition, err := t.ts.GetCondition(tenderId, store.Latest)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return "", services.ErrServiceDatabaseDisconnected
		}
		if errors.Is(err, store.ErrRecordNotFound) {
			return "", services.ErrNoSuchTender
		}
		t.logger.Errorf("unexpected error: %s on method GetCondition", err)
		return "", err
	}

	if tenderCondition.Status != "Published" && tenderCondition.OrgId != orgId {
		return "", services.ErrNoPermitions
	}

	return tenderCondition.Status, nil
}

func (t *Tender) ChangeStat(tenderId, status, username string) (*models.Tender, error) {
	orgId, err := t.rs.GetOrgId(username)
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
		t.logger.Errorf("unexpected error: %s on method GetOrgId", err)
		return nil, err
	}

	tenderCondition, err := t.ts.GetCondition(tenderId, store.Latest)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		if errors.Is(err, store.ErrRecordNotFound) {
			return nil, services.ErrNoSuchTender
		}
		t.logger.Errorf("unexpected error: %s on method GetCondition", err)
		return nil, err
	}

	if tenderCondition.Status == status {
		return tenderCondition, nil
	}

	if tenderCondition.OrgId != orgId {
		return nil, services.ErrNoPermitions
	}

	tenderCondition.Status = status
	tenderCondition.Version += 1

	result, err := t.ts.UpdateCondition(tenderCondition)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		t.logger.Errorf("unexpected error: %s on method UpdateCondition", err)
		return nil, err
	}
	return result, nil
}

func (t *Tender) Edit(tnd *models.Tender, tenderId, username string) (*models.Tender, error) {
	orgId, err := t.rs.GetOrgId(username)
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
		t.logger.Errorf("unexpected error: %s on method GetOrgId", err)
		return nil, err
	}

	tenderCondition, err := t.ts.GetCondition(tenderId, store.Latest)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		if errors.Is(err, store.ErrRecordNotFound) {
			return nil, services.ErrNoSuchTender
		}
		t.logger.Errorf("unexpected error: %s on method GetCondition", err)
		return nil, err
	}

	if tenderCondition.OrgId != orgId {
		return nil, services.ErrNoPermitions
	}

	var count int
	if tnd.Name != "" && tenderCondition.Name != tnd.Name {
		tenderCondition.Name = tnd.Name
		count++
	}
	if tnd.Description != "" && tenderCondition.Description != tnd.Description {
		tenderCondition.Description = tnd.Description
		count++
	}
	if tnd.ServType != "" && tenderCondition.ServType != tnd.ServType {
		tenderCondition.ServType = tnd.ServType
		count++
	}
	if count == 0 {
		return tenderCondition, nil
	}

	tenderCondition.Version += 1

	result, err := t.ts.UpdateCondition(tenderCondition)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		t.logger.Errorf("unexpected error: %s on method UpdateCondition", err)
		return nil, err
	}
	return result, nil
}

func (t *Tender) Rollback(tenderId, username string, version int64) (*models.Tender, error) {
	orgId, err := t.rs.GetOrgId(username)
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
		t.logger.Errorf("unexpected error: %s on method GetOrgId", err)
		return nil, err
	}

	tenderCondition, err := t.ts.GetCondition(tenderId, version)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		if errors.Is(err, store.ErrRecordNotFound) {
			return nil, services.ErrNoSuchTender
		}
		t.logger.Errorf("unexpected error: %s on method GetCondition", err)
		return nil, err
	}

	if tenderCondition.OrgId != orgId {
		return nil, services.ErrNoPermitions
	}

	latestVersion, err := t.ts.GetTenderLatestVersion(tenderId, username)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		t.logger.Errorf("unexpected error: %s on method GetTenderLatestVersion", err)
		return nil, err
	}

	tenderCondition.Version = latestVersion + 1

	result, err := t.ts.UpdateCondition(tenderCondition)
	if err != nil {
		if errors.Is(err, store.ErrConnClosed) {
			return nil, services.ErrServiceDatabaseDisconnected
		}
		t.logger.Errorf("unexpected error: %s on method UpdateCondition", err)
		return nil, err
	}
	return result, nil
}
