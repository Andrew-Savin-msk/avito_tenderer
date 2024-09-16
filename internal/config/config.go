package config

import (
	"log"
	"os"
	"strings"
)

type Server struct {
	Port string
}

type Database struct {
	Conn string
}

type Config struct {
	Srv Server
	Db  Database
}

func Load() *Config {
	srvAddr := getEnv("SERVER_ADDRESS")
	div := strings.Index(srvAddr, ":")
	if div == -1 {
		log.Fatal("incorrect server address")
	}

	config := &Config{
		Srv: Server{
			Port: srvAddr[div+1:],
		},
		Db: Database{
			Conn: getEnv("POSTGRES_CONN"),
		},
	}

	return config
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatal("error: ENV variable not found")
	}
	return value
}
