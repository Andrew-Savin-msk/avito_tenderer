package responsiblestore

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/domain/models"
	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/store"
)

type ResponsibleStore struct {
	db *sql.DB
}

func New(db *sql.DB) *ResponsibleStore {
	return &ResponsibleStore{
		db: db,
	}
}

func (r *ResponsibleStore) GetResponsibleUUID(responsibles *models.Responsible) (string, error) {
	var userId string
	err := r.db.QueryRow(
		"SELECT id FROM employee WHERE username = $1;",
		responsibles.Username,
	).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", store.ErrUserNotFound
		}
		if errors.Is(err, sql.ErrConnDone) {
			return "", store.ErrConnClosed
		}
		return "", err
	}

	fmt.Println("user found")

	var resposiblesUUID string
	err = r.db.QueryRow(
		"SELECT id FROM organization_responsible WHERE user_id = $1 AND organization_id = $2;",
		userId,
		responsibles.OrgId,
	).Scan(&resposiblesUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", store.ErrRecordNotFound
		}
		if errors.Is(err, sql.ErrConnDone) {
			return "", store.ErrConnClosed
		}
		return "", err
	}
	return resposiblesUUID, nil
}

func (r *ResponsibleStore) GetRespUUIDs(username string) ([]string, error) {
	var userId string
	err := r.db.QueryRow(
		"SELECT id FROM employee WHERE username = $1;",
		username,
	).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrUserNotFound
		}
		if errors.Is(err, sql.ErrConnDone) {
			return nil, store.ErrConnClosed
		}
		return nil, err
	}

	rows, err := r.db.Query(
		"SELECT id FROM organization_responsible WHERE user_id = $1;",
		userId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrRecordNotFound
		}
		if errors.Is(err, sql.ErrConnDone) {
			return nil, store.ErrConnClosed
		}
		return nil, err
	}

	var resposiblesUUIDs []string
	for rows.Next() {
		var uuid string
		err = rows.Scan(&uuid)
		if err != nil {
			return nil, err
		}
		resposiblesUUIDs = append(resposiblesUUIDs, uuid)
	}

	return resposiblesUUIDs, nil
}

func (r *ResponsibleStore) IsResponcible(responsible *models.Responsible) error {
	var userId string
	err := r.db.QueryRow(
		"SELECT id FROM employee WHERE username = $1;",
		responsible.Username,
	).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return store.ErrUserNotFound
		}
		if errors.Is(err, sql.ErrConnDone) {
			return store.ErrConnClosed
		}
		return err
	}

	var respId string
	err = r.db.QueryRow(
		"SELECT id FROM organization_responsible WHERE user_id = $1 AND organization_id = $2;",
		userId,
		responsible.OrgId,
	).Scan(&respId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return store.ErrRecordNotFound
		}
		if errors.Is(err, sql.ErrConnDone) {
			return store.ErrConnClosed
		}
		return err
	}
	return nil
}

func (r *ResponsibleStore) IsUserExists(username string) error {
	var userId string
	err := r.db.QueryRow(
		"SELECT id FROM employee WHERE username = $1;",
		username,
	).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return store.ErrUserNotFound
		}
		if errors.Is(err, sql.ErrConnDone) {
			return store.ErrConnClosed
		}
		return err
	}
	return nil
}

func (r *ResponsibleStore) GetOrgId(username string) (string, error) {
	userId, err := r.GetUserId(username)
	if err != nil {
		return "", err
	}

	return r.ResponcibleForOrg(userId)
}

func (r *ResponsibleStore) ResponcibleForOrg(userId string) (string, error) {
	var orgId string
	err := r.db.QueryRow(
		"SELECT organization_id FROM organization_responsible WHERE user_id = $1;",
		userId,
	).Scan(&orgId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", store.ErrRecordNotFound
		}
		if errors.Is(err, sql.ErrConnDone) {
			return "", store.ErrConnClosed
		}
		return "", err
	}
	return orgId, nil
}

func (r *ResponsibleStore) GetUserId(username string) (string, error) {
	var userId string
	err := r.db.QueryRow(
		"SELECT id FROM employee WHERE username = $1;",
		username,
	).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", store.ErrUserNotFound
		}
		if errors.Is(err, sql.ErrConnDone) {
			return "", store.ErrConnClosed
		}
		return "", err
	}

	return userId, nil
}
