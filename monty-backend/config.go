package main

import (
	"crypto/rand"
	"os"
	"strconv"
)

type Config struct {
	Port      int
	StaticDir string
	DB        DBConfig
	JWTSecret []byte
}

type DBConfig struct {
	User     string
	Password string
	Database string
	Port     int
	MaxConns int
	Host     string
}

type Context struct {
	Config *Config
	DB     *DB
}

const (
	DEFAULT_PORT = 8080
	STATIC_DIR   = "./static"

	DEFAULT_DB_USER     = "postgres"
	DEFAULT_DB_DATABASE = "db"
	DEFAULT_DB_HOST     = "localhost"
	MAX_DB_CONNECTIONS  = 5
	DEFAULT_DB_PORT     = 5432
)

func configParseInt(name string, def int) int {
	if val, ok := os.LookupEnv(name); ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return def
}

func configParseString(name string, def string) string {
	if val, ok := os.LookupEnv(name); ok {
		return val
	}
	return def
}

func NewConfigFromEnv() *Config {
	var jwtSecretBytes []byte

	port := configParseInt("PORT", DEFAULT_PORT)
	dbPort := configParseInt("DB_PORT", DEFAULT_DB_PORT)
	dbUser := configParseString("DB_USER", DEFAULT_DB_USER)
	staticDir := configParseString("STATIC_DIR", STATIC_DIR)
	dbDatabase := configParseString("DB_DATABASE", DEFAULT_DB_DATABASE)
	dbMaxConn := configParseInt("DB_MAX_CONNS", MAX_DB_CONNECTIONS)
	dbHost := configParseString("DB_HOST", DEFAULT_DB_HOST)

	dbPassword, _ := os.LookupEnv("DB_PASSWORD")

	if dbPassword == "" {
		dbPassword = "password"
	}

	jwtSecret, _ := os.LookupEnv("JWT_SECRET")

	if jwtSecret == "" {
		jwtSecretBytes = make([]byte, 32)
		rand.Read(jwtSecretBytes)
	} else {
		jwtSecretBytes = []byte(jwtSecret)
	}

	return &Config{
		Port:      port,
		StaticDir: staticDir,
		JWTSecret: jwtSecretBytes,
		DB: DBConfig{
			User:     dbUser,
			Password: dbPassword,
			Database: dbDatabase,
			Port:     dbPort,
			MaxConns: dbMaxConn,
			Host:     dbHost,
		},
	}
}
