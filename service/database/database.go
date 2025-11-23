package database

import (
	"database/sql"
	"fmt"
)

// User rappresenta la riga della tabella 'users'
type User struct {
	ID       string
	Name     string
	PhotoURL string
}

// AppDatabase definisce l'interfaccia del database
type AppDatabase interface {
	Ping() error

	// --- USERS ---
	CreateUser(u User) (User, error)
	GetUserByName(name string) (User, error)
	GetUserByID(id string) (User, error)
	SetUserName(id string, name string) error
	SetPhotoUrl(id string, url string) error
}

type appdbimpl struct {
	c *sql.DB
}

// New crea una nuova istanza di AppDatabase accettando una connessione già aperta
// NOTA: La firma è cambiata da (string) a (*sql.DB)
func New(db *sql.DB) (AppDatabase, error) {

	// 1. Verifica se la connessione passata è viva
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// 2. Abilita Foreign Keys (Essenziale per SQLite)
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, fmt.Errorf("error enabling foreign keys: %w", err)
	}

	// 3. Creazione Tabelle (Schema Migration)
	var tableName = "users"
	var tableQuery = `
		CREATE TABLE IF NOT EXISTS users (
			id TEXT NOT NULL PRIMARY KEY,
			name TEXT NOT NULL,
			photo_url TEXT
		);
	`

	if _, err := db.Exec(tableQuery); err != nil {
		return nil, fmt.Errorf("error creating table %s: %w", tableName, err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
