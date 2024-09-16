package tenderstore

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/domain/models"
	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/store"
	"github.com/lib/pq"
)

type TenderStore struct {
	db    *sql.DB
	stats map[string]string
}

func New(db *sql.DB) *TenderStore {
	return &TenderStore{
		db: db,
		stats: map[string]string{
			"CREATED":   "Created",
			"PUBLISHED": "Published",
			"CLOSED":    "Closed",
		},
	}
}

func (t *TenderStore) GetLimitedList(limit, offset int64, servType []string) ([]*models.Tender, error) {
	var err error
	var rows *sql.Rows
	if len(servType) != 0 {
		rows, err = t.db.Query(
			"SELECT tv.tender_id, tv.name, tv.description, tv.status, tv.type, tv.version, tv.created_at "+
				"FROM tenders_versions AS tv "+
				"INNER JOIN ( "+
				"SELECT tender_id, MAX(version) AS latest_version "+
				"FROM tenders_versions "+
				"GROUP BY tender_id"+
				") AS lv "+
				"ON tv.tender_id = lv.tender_id AND tv.version = lv.latest_version "+
				"WHERE tv.type = ANY($1) "+
				"ORDER BY tv.name ASC "+
				"LIMIT $2 "+
				"OFFSET $3;",
			pq.Array(servType),
			limit,
			offset,
		)
	} else {
		rows, err = t.db.Query(
			"SELECT tv.tender_id, tv.name, tv.description, tv.status, tv.type, tv.version, tv.created_at "+
				"FROM tenders_versions AS tv "+
				"INNER JOIN ( "+
				"SELECT tender_id, MAX(version) AS latest_version "+
				"FROM tenders_versions "+
				"GROUP BY tender_id"+
				") AS lv "+
				"ON tv.tender_id = lv.tender_id AND tv.version = lv.latest_version "+
				"ORDER BY tv.name ASC "+
				"LIMIT $1 "+
				"OFFSET $2;",
			limit,
			offset,
		)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrRecordNotFound
		}
		if strings.Contains(err.Error(), "no such host") {
			return nil, store.ErrConnClosed
		}
		return nil, err
	}

	result := []*models.Tender{}
	for rows.Next() {
		var tender models.Tender
		err = rows.Scan(&tender.Id, &tender.Name, &tender.Description, &tender.Status, &tender.ServType, &tender.Version, &tender.Created)
		if err != nil {
			return nil, err
		}
		tender.Status = t.stats[tender.Status]
		result = append(result, &tender)
	}
	return result, nil
}

func (t *TenderStore) Create(tnd *models.Tender, resp *models.Responsible) (*models.Tender, error) {
	err := t.db.QueryRow(
		"INSERT INTO tenders (organization_id, username) VALUES ($1, $2) RETURNING id;",
		resp.OrgId,
		resp.Username,
	).Scan(&tnd.Id)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return nil, store.ErrConnClosed
		}
		return nil, err
	}

	err = t.db.QueryRow(
		"INSERT INTO tenders_versions (tender_id, name, description, status, type) VALUES ($1, $2, $3, 'CREATED', $4) RETURNING created_at, version, status;",
		tnd.Id,
		tnd.Name,
		tnd.Description,
		tnd.ServType,
	).Scan(&tnd.Created, &tnd.Version, &tnd.Status)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return nil, store.ErrConnClosed
		}
		return nil, err
	}
	tnd.Status = t.stats[tnd.Status]

	return tnd, nil
}

func (t *TenderStore) GetUserTenders(limit, offset int64, username string) ([]*models.Tender, error) {
	rows, err := t.db.Query(
		"SELECT tv.tender_id, tv.name, tv.description, tv.status, tv.type, tv.version, tv.created_at "+
			"FROM tenders AS t "+
			"INNER JOIN tenders_versions AS tv "+
			"ON t.id = tv.tender_id "+
			"INNER JOIN ( "+
			"SELECT tender_id, MAX(version) AS latest_version "+
			"FROM tenders_versions "+
			"GROUP BY tender_id"+
			") AS lv "+
			"ON tv.tender_id = lv.tender_id AND tv.version = lv.latest_version "+
			"WHERE t.username = $1 "+
			"ORDER BY tv.name ASC "+
			"LIMIT $2 "+
			"OFFSET $3;",
		username,
		limit,
		offset,
	)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return nil, store.ErrConnClosed
		}
		return nil, err
	}

	result := []*models.Tender{}
	for rows.Next() {
		var tender models.Tender
		err = rows.Scan(&tender.Id, &tender.Name, &tender.Description, &tender.Status, &tender.ServType, &tender.Version, &tender.Created)
		if err != nil {
			return nil, err
		}
		tender.Status = t.stats[tender.Status]
		result = append(result, &tender)
	}
	return result, nil
}

func (t *TenderStore) GetStatus(tenderId string) (string, error) {
	var status string
	err := t.db.QueryRow(
		"SELECT tv.status "+
			"FROM tenders_versions AS tv "+
			"JOIN ( "+
			"SELECT tender_id, MAX(version) AS latest_version "+
			"FROM tenders_versions "+
			"GROUP BY tender_id"+
			") AS lv "+
			"ON tv.tender_id = lv.tender_id AND tv.version = lv.latest_version "+
			"WHERE tv.tender_id = $1;",
		tenderId,
	).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", store.ErrRecordNotFound
		}
		if strings.Contains(err.Error(), "no such host") {
			return "", store.ErrConnClosed
		}
		return "", err
	}
	return t.stats[status], nil
}

func (t *TenderStore) IsResponcibleFor(tenderId string, respUUIDs []string) error {
	if len(respUUIDs) == 0 {
		return store.ErrRecordNotFound
	}
	err := t.db.QueryRow(
		"SELECT id FROM tenders WHERE id = $1 AND responsible_id = ANY($2)",
		tenderId,
		pq.Array(respUUIDs),
	).Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return store.ErrRecordNotFound
		}
		if strings.Contains(err.Error(), "no such host") {
			return store.ErrConnClosed
		}
		return err
	}
	return nil
}

func (t *TenderStore) GetTenderLatestVersion(tenderId, username string) (int64, error) {
	var version int64
	err := t.db.QueryRow(
		"SELECT tv.version "+
			"FROM tenders AS t "+
			"INNER JOIN tenders_versions tv ON t.id = tv.tender_id "+
			"JOIN ( "+
			"SELECT tender_id, MAX(version) AS latest_version "+
			"FROM tenders_versions "+
			"GROUP BY tender_id"+
			") lv ON tv.tender_id = lv.tender_id AND tv.version = lv.latest_version "+
			"WHERE t.id = $1 AND t.username = $2;",
		tenderId,
		username,
	).Scan(&version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, store.ErrRecordNotFound
		}
		if strings.Contains(err.Error(), "no such host") {
			return -1, store.ErrConnClosed
		}
		return -1, err
	}
	return version, nil
}
func (t *TenderStore) GetCondition(tenderId string, version int64) (*models.Tender, error) {
	var tnd models.Tender
	var err error
	if version == store.Latest {
		err = t.db.QueryRow(
			"SELECT tv.tender_id, tv.name, tv.description, tv.status, tv.type, t.organization_id, tv.version, tv.created_at "+
				"FROM tenders AS t "+
				"INNER JOIN tenders_versions tv ON t.id = tv.tender_id "+
				"JOIN ( "+
				"SELECT tender_id, MAX(version) AS latest_version "+
				"FROM tenders_versions "+
				"GROUP BY tender_id"+
				") lv ON tv.tender_id = lv.tender_id AND tv.version = lv.latest_version "+
				"WHERE t.id = $1;",
			tenderId,
		).Scan(&tnd.Id, &tnd.Name, &tnd.Description, &tnd.Status, &tnd.ServType, &tnd.OrgId, &tnd.Version, &tnd.Created)
	} else {
		err = t.db.QueryRow(
			"SELECT tv.tender_id, tv.name, tv.description, tv.status, tv.type, t.organization_id, tv.version, tv.created_at "+
				"FROM tenders AS t "+
				"INNER JOIN tenders_versions tv ON t.id = tv.tender_id "+
				"WHERE t.id = $1 AND tv.version = $2;",
			tenderId,
			version,
		).Scan(&tnd.Id, &tnd.Name, &tnd.Description, &tnd.Status, &tnd.ServType, &tnd.OrgId, &tnd.Version, &tnd.Created)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrRecordNotFound
		}
		if strings.Contains(err.Error(), "no such host") {
			return nil, store.ErrConnClosed
		}
		return nil, err
	}
	tnd.Status = t.stats[tnd.Status]
	return &tnd, nil
}

// GetResponsible(tenderId string) (*models.Responsible, error)
func (t *TenderStore) UpdateCondition(newCondition *models.Tender) (*models.Tender, error) {
	newCondition.Status = strings.ToUpper(newCondition.Status)

	_, err := t.db.Exec(
		"INSERT INTO tenders_versions (tender_id, name, description, status, type, version, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		newCondition.Id,
		newCondition.Name,
		newCondition.Description,
		newCondition.Status,
		newCondition.ServType,
		newCondition.Version,
		newCondition.Created,
	)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return nil, store.ErrConnClosed
		}
		return nil, err
	}
	newCondition.Status = t.stats[newCondition.Status]
	return newCondition, nil
}

func (t *TenderStore) GetOrgIdByBidId(bidId string) (string, error) {
	var orgId string
	err := t.db.QueryRow(
		"SELECT t.organization_id "+
			"FROM bids AS b "+
			"INNER JOIN tenders t ON b.tender_id = t.id "+
			"WHERE b.id = $1;",
		bidId,
	).Scan(&orgId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", store.ErrRecordNotFound
		}
		if strings.Contains(err.Error(), "no such host") {
			return "", store.ErrConnClosed
		}
		return "", err
	}
	return orgId, nil
}
