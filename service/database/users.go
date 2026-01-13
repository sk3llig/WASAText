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

// CheckUserExists verifica se un utente esiste dato il suo ID
func (db *appdbimpl) CheckUserExists(id string) (bool, error) {
	var exists int
	// La query cerca 1 riga con quell'ID
	err := db.c.QueryRow("SELECT 1 FROM users WHERE id = ?", id).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil // Non c'è errore, ma non esiste l'utente
	}
	if err != nil {
		return false, err // Errore vero del database
	}
	return true, nil // Esiste!
}

// CheckUserByUsername verifica se esiste un utente dato il suo USERNAME
func (db *appdbimpl) CheckUserByUsername(username string) (bool, error) {
	var exists int
	// MODIFICA QUI: cerchiamo nella colonna 'username', non 'id'
	err := db.c.QueryRow("SELECT 1 FROM users WHERE name = ?", username).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil // Username non trovato
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
