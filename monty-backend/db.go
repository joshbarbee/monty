package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

const (
	CREATE_SCHEMA = `CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password BYTEA NOT NULL,
		admin BOOLEAN DEFAULT FALSE,
		salt BYTEA NOT NULL,
		status INT DEFAULT 0,
		created TIMESTAMPTZ,
		modified TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		accessed TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	)`

	UPDATE_LAST_ACCESSED = "UPDATE users SET accessed = CURRENT_TIMESTAMP WHERE id = $1"

	CREATE_USER = `INSERT INTO users (id, username, email, password, salt, created) 
					VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
)

type DB struct {
	Conn *pgx.ConnPool
}

func NewDB(config DBConfig) (*DB, error) {
	dbConfig := pgx.ConnConfig{
		User:     config.User,
		Password: config.Password,
		Database: config.Database,
		Host:     config.Host,
		Port:     uint16(config.Port),
	}

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     dbConfig,
		MaxConnections: config.MaxConns,
	}

	pool, err := pgx.NewConnPool(poolConfig)

	if err != nil {
		return nil, err
	}

	return &DB{pool}, nil
}

func (db *DB) Close() {
	db.Conn.Close()
}

func (db *DB) CreateSchema() error {
	_, err := db.Conn.Exec(CREATE_SCHEMA)

	return err
}

func (db *DB) UpdateLastAccessed(id uuid.UUID) error {
	_, err := db.Conn.Exec(UPDATE_LAST_ACCESSED, id)

	return err
}

func (db *DB) UserExists(username string) error {
	err := db.Conn.QueryRow("SELECT username FROM users WHERE username = $1", username).Scan(&username)

	return err
}

func (db *DB) CreateUser(user *User) error {
	t := time.Now()

	err := db.UserExists(user.Username)
	if err == nil {
		return &ErrUserExists{user.Username}
	}

	_, err = db.Conn.Exec(CREATE_USER, user.ID, user.Username, user.Email, user.Password, user.Salt, t)

	user.Created = t

	return err
}

func (db *DB) VerifyUser(username, password string) (User, error) {
	var user User
	t := time.Now()

	err := db.Conn.QueryRow("SELECT id, username, email, password, salt FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Salt)

	if err != nil {
		return user, &ErrInvalidPassword{}
	}

	if !user.VerifyPassword(password) {
		return user, &ErrInvalidPassword{}
	}

	err = db.UpdateLastAccessed(user.ID)
	user.Accessed = t

	return user, err
}
