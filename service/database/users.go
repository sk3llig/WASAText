package database

import (
	"database/sql"
	"errors"
)

// CreateUser inserisce un nuovo utente nel database
func (db *appdbimpl) CreateUser(u User) (User, error) {
	// La query INSERT
	query := `INSERT INTO users (id, name, photo_url) VALUES (?, ?, ?)`

	// Esecuzione (sostituisce i ? con i valori)
	_, err := db.c.Exec(query, u.ID, u.Name, u.PhotoURL)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

// GetUserByName cerca un utente dato il suo nome (per il login)
func (db *appdbimpl) GetUserByName(name string) (User, error) {
	var u User

	// Query SELECT
	query := `SELECT id, name, photo_url FROM users WHERE name = ?`

	// QueryRow esegue la query e si aspetta una sola riga
	err := db.c.QueryRow(query, name).Scan(&u.ID, &u.Name, &u.PhotoURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Importante: ritorniamo un errore specifico se non esiste
			return User{}, errors.New("user not found")
		}
		return User{}, err
	}

	return u, nil
}

// GetUserByID cerca un utente dato il suo ID (per i messaggi/profilo)
func (db *appdbimpl) GetUserByID(id string) (User, error) {
	var u User
	query := `SELECT id, name, photo_url FROM users WHERE id = ?`

	err := db.c.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.PhotoURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, errors.New("user not found")
		}
		return User{}, err
	}
	return u, nil
}

// SetUserName aggiorna il nome utente
func (db *appdbimpl) SetUserName(id string, name string) error {
	query := `UPDATE users SET name = ? WHERE id = ?`
	_, err := db.c.Exec(query, name, id)
	return err
}

// SetPhotoUrl aggiorna la foto profilo
func (db *appdbimpl) SetPhotoUrl(id string, url string) error {
	query := `UPDATE users SET photo_url = ? WHERE id = ?`
	_, err := db.c.Exec(query, url, id)
	return err
}
