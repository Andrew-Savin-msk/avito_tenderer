package bidstore

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/domain/models"
	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/store"
)

type BidStore struct {
	db    *sql.DB
	stats map[string]string
}

func New(db *sql.DB) *BidStore {
	return &BidStore{
		db: db,
		stats: map[string]string{
			"CREATED":   "Created",
			"PUBLISHED": "Published",
			"CANCELED":  "Canceled",
		},
	}
}

func (b *BidStore) Create(bid *models.Bid, orgId string) (*models.Bid, error) {
	var err error
	if bid.AuthorType == "Organization" {
		err = b.db.QueryRow(
			"INSERT INTO bids (tender_id, author_type, organization_id) VALUES ($1, $2, $3) RETURNING id;",
			bid.TenderId,
			bid.AuthorType,
			bid.AuthorId,
		).Scan(&bid.Id)
		if err != nil {
			if strings.Contains(err.Error(), "no such host") {
				return nil, store.ErrConnClosed
			}
			return nil, err
		}
	} else {
		err = b.db.QueryRow(
			"INSERT INTO bids (tender_id, author_type, organization_id, user_id) VALUES ($1, $2, $3, $4) RETURNING id;",
			bid.TenderId,
			bid.AuthorType,
			orgId,
			bid.AuthorId,
		).Scan(&bid.Id)
		if err != nil {
			if strings.Contains(err.Error(), "no such host") {
				return nil, store.ErrConnClosed
			}
			return nil, err
		}
	}

	err = b.db.QueryRow(
		"INSERT INTO bids_versions (bid_id, name, description, status) VALUES ($1, $2, $3, 'CREATED') RETURNING created_at, version, status;",
		bid.Id,
		bid.Name,
		bid.Description,
	).Scan(&bid.Created, &bid.Version, &bid.Status)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return nil, store.ErrConnClosed
		}
		return nil, err
	}
	bid.Status = b.stats[bid.Status]

	return bid, nil
}

func (b *BidStore) GetUserList(limit, offset int64, userId string) ([]*models.Bid, error) {
	rows, err := b.db.Query(
		"SELECT bv.bid_id, bv.name, bv.description, bv.status, b.author_type, b.user_id, bv.version, bv.created_at "+
			"FROM bids_versions bv "+
			"INNER JOIN ( "+
			"SELECT bid_id, MAX(version) AS latest_version "+
			"FROM bids_versions "+
			"GROUP BY bid_id"+
			") lv ON bv.bid_id = lv.bid_id AND bv.version = lv.latest_version "+
			"INNER JOIN bids b ON b.id = bv.bid_id "+
			"WHERE b.author_type = 'User' AND b.user_id = $1 "+
			"ORDER BY bv.name ASC "+
			"LIMIT $2 "+
			"OFFSET $3;",
		userId,
		limit,
		offset,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrRecordNotFound
		}
		if strings.Contains(err.Error(), "no such host") {
			return nil, store.ErrConnClosed
		}
		return nil, err
	}

	result := []*models.Bid{}
	for rows.Next() {
		var bid models.Bid
		err = rows.Scan(&bid.Id, &bid.Name, &bid.Description, &bid.Status, &bid.AuthorType, &bid.AuthorId, &bid.Version, &bid.Created)
		if err != nil {
			return nil, err
		}
		bid.Status = b.stats[bid.Status]
		result = append(result, &bid)
	}
	return result, nil
}

func (b *BidStore) GetCondition(bidId string, version int64) (*models.Bid, error) {
	var err error
	var bid models.Bid
	if version == store.Latest {
		err = b.db.QueryRow(
			"SELECT bv.bid_id, b.tender_id, bv.name, bv.description, bv.status, b.author_type, CASE WHEN b.author_type = 'User' THEN b.user_id ELSE b.organization_id END AS author_id, bv.version, bv.created_at "+
				"FROM bids_versions bv "+
				"INNER JOIN ( "+
				"SELECT bid_id, MAX(version) AS latest_version "+
				"FROM bids_versions "+
				"GROUP BY bid_id"+
				") lv ON bv.bid_id = lv.bid_id AND bv.version = lv.latest_version "+
				"INNER JOIN bids b ON b.id = bv.bid_id "+
				"WHERE bv.bid_id = $1;",
			bidId,
		).Scan(&bid.Id, &bid.TenderId, &bid.Name, &bid.Description, &bid.Status, &bid.AuthorType, &bid.AuthorId, &bid.Version, &bid.Created)
	} else {
		err = b.db.QueryRow(
			"SELECT bv.bid_id, b.tender_id, bv.name, bv.description, bv.status, b.author_type, CASE WHEN b.author_type = 'User' THEN b.user_id ELSE b.organization_id END AS author_id, bv.version, bv.created_at "+
				"FROM bids_versions bv "+
				"INNER JOIN bids b ON b.id = bv.bid_id "+
				"WHERE b.bid_id = $1 AND bv.version = $2;",
			bidId,
			version,
		).Scan(&bid.Id, &bid.TenderId, &bid.Name, &bid.Description, &bid.Status, &bid.AuthorType, &bid.AuthorId, &bid.Version, &bid.Created)
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
	bid.Status = b.stats[bid.Status]
	return &bid, nil
}

func (b *BidStore) GetBidLatestVersion(bidId string) (int64, error) {
	var version int64
	err := b.db.QueryRow(
		"SELECT bv.version "+
			"FROM bids_versions bv "+
			"INNER JOIN ( "+
			"SELECT bid_id, MAX(version) AS latest_version "+
			"FROM bids_versions "+
			"GROUP BY bid_id"+
			") lv ON bv.bid_id = lv.bid_id AND bv.version = lv.latest_version "+
			"INNER JOIN bids b ON b.id = bv.bid_id "+
			"WHERE bv.bid_id = $1;",
		bidId,
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

func (b *BidStore) UpdateCondition(newCondition *models.Bid) (*models.Bid, error) {
	newCondition.Status = strings.ToUpper(newCondition.Status)

	_, err := b.db.Exec(
		"INSERT INTO bids_versions (bid_id, name, description, status, version, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		newCondition.Id,
		newCondition.Name,
		newCondition.Description,
		newCondition.Status,
		newCondition.Version,
		newCondition.Created,
	)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return nil, store.ErrConnClosed
		}
		return nil, err
	}
	newCondition.Status = b.stats[newCondition.Status]
	return newCondition, nil
}

func (b *BidStore) GetTenderList(limit, offset int64, tenderId, orgId string) ([]*models.Bid, error) {
	rows, err := b.db.Query(
		"SELECT bv.bid_id, bv.name, bv.description, bv.status, b.author_type, CASE WHEN b.author_type = 'User' THEN b.user_id ELSE b.organization_id END AS author_id, bv.version, bv.created_at "+
			"FROM bids_versions AS bv "+
			"INNER JOIN ( "+
			"SELECT bid_id, MAX(version) AS latest_version "+
			"FROM bids_versions "+
			"GROUP BY bid_id "+
			") lv ON bv.bid_id = lv.bid_id AND bv.version = lv.latest_version "+
			"INNER JOIN bids AS b ON b.id = bv.bid_id "+
			"WHERE b.tender_id = $1 AND (bv.status = 'PUBLISHED' OR b.organization_id = $2) "+
			"ORDER BY bv.name ASC "+
			"LIMIT $3 "+
			"OFFSET $4",
		tenderId,
		orgId,
		limit,
		offset,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrRecordNotFound
		}
		if strings.Contains(err.Error(), "no such host") {
			return nil, store.ErrConnClosed
		}
		return nil, err
	}

	result := []*models.Bid{}
	for rows.Next() {
		var bid models.Bid
		err = rows.Scan(&bid.Id, &bid.Name, &bid.Description, &bid.Status, &bid.AuthorType, &bid.AuthorId, &bid.Version, &bid.Created)
		if err != nil {
			return nil, err
		}
		bid.Status = b.stats[bid.Status]
		result = append(result, &bid)
	}
	return result, nil
}

func (b *BidStore) AddFeedback(bidId, userId, feedback string) error {
	_, err := b.db.Exec(
		"INSERT INTO feedbacks (bid_id, user_id, feedback) VALUES ($1, $2, $3);",
		bidId,
		userId,
		feedback,
	)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return store.ErrConnClosed
		}
		return err
	}
	return nil
}

func (b *BidStore) GetFeedbacks(tenderId, authorUserId string, limit, offset int64) ([]*models.Feedback, error) {
	rows, err := b.db.Query(
		"SELECT f.id, f.feedback, f.created_at "+
			"FROM feedbacks AS f "+
			"INNER JOIN bids AS b ON f.bid_id = b.id "+
			"INNER JOIN bids_versions bv ON b.id = bv.bid_id "+
			"INNER JOIN tenders AS t ON b.tender_id = t.id "+
			"WHERE bv.status = 'PUBLISHED' AND t.id = $1 "+
			"ORDER BY f.feedback ASC "+
			"LIMIT $2 "+
			"OFFSET $3;",
		tenderId,
		limit,
		offset,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrRecordNotFound
		}
		if strings.Contains(err.Error(), "no such host") {
			return nil, store.ErrConnClosed
		}
		return nil, err
	}

	result := []*models.Feedback{}
	for rows.Next() {
		var feed models.Feedback
		err = rows.Scan(&feed.Id, &feed.Desc, &feed.Created)
		if err != nil {
			return nil, err
		}
		result = append(result, &feed)
	}
	return result, nil
}
