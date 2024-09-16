package apiserver

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/config"
	bidservice "github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/services/bider"
	tenderservice "github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/services/tender"
	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/store"
	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/store/bidstore"
	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/store/responsiblestore"
	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/store/tenderstore"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// Start init's all connections and starts api's work
func Start(cfg *config.Config) error {
	// Get logger
	log := setLog("debug")

	// Get db connection
	tenderSt, err := loadTenderStore(cfg.Db)
	if err != nil {
		return fmt.Errorf("unable to load tenders store error: %s", err)
	}

	responsibleSt, err := loadResponsibleStore(cfg.Db)
	if err != nil {
		return fmt.Errorf("unable to load responsibles store error: %s", err)
	}

	bidSt, err := loadBidStore(cfg.Db)
	if err != nil {
		return fmt.Errorf("unable to load responsibles store error: %s", err)
	}

	// Get Tender Service
	TenderServ := tenderservice.New(tenderSt, responsibleSt, log)
	BidsServ := bidservice.New(tenderSt, bidSt, responsibleSt, log)

	// Get Bid Service

	// Get server
	srv := newServer(log, TenderServ, BidsServ)

	log.Infof("api strted work on port: %s", cfg.Srv.Port)

	// Start listner
	err = http.ListenAndServe(":"+cfg.Srv.Port, srv)
	if err != nil {
		log.Infof("api ended work with error: %s", err)
	} else {
		log.Info("api ended work")
	}

	return nil
}

func setLog(level string) *logrus.Logger {
	log := logrus.New()
	switch strings.ToLower(level) {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	}
	fmt.Printf("logger set in level: %s\n", level)
	return log
}

func loadResponsibleStore(cfg config.Database) (store.Responsibles, error) {
	db, err := sql.Open("postgres", cfg.Conn)
	if err != nil {
		return nil, fmt.Errorf("open: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return responsiblestore.New(db), nil
}

func loadTenderStore(cfg config.Database) (store.Tenders, error) {
	db, err := sql.Open("postgres", cfg.Conn)
	if err != nil {
		return nil, fmt.Errorf("open: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return tenderstore.New(db), nil
}

func loadBidStore(cfg config.Database) (store.Bids, error) {
	db, err := sql.Open("postgres", cfg.Conn)
	if err != nil {
		return nil, fmt.Errorf("open: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return bidstore.New(db), nil
}
