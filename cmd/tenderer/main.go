package main

import (
	"log"

	apiserver "github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/api_server"
	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/config"
)

func main() {
	cfg := config.Load()

	err := apiserver.Start(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
