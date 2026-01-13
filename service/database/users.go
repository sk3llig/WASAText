package database

import (
	"database/sql"
	"errors"
	"fmt"
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
func (db *appdbimpl) SetUserName(userID string, newName string) error {

	// CORREZIONE: Qui c'era scritto quasi sicuramente "username"
	// Devi cambiarlo in "name" perché è così che si chiama la tua colonna nel DB
	query := `UPDATE users SET name = ? WHERE id = ?`

	_, err := db.c.Exec(query, newName, userID)
	if err != nil {
		return fmt.Errorf("error updating name: %w", err)
	}
	return nil
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

// GetUserIDByName cerca l'ID dato il nome.
// Se non lo trova, restituisce stringa vuota e NESSUN errore.
func (db *appdbimpl) GetUserIDByName(name string) (string, error) {
	var id string
	// Assicurati che la colonna sia 'name'
	query := "SELECT id FROM users WHERE name = ?"

	err := db.c.QueryRow(query, name).Scan(&id)

	// --- QUESTA È LA PARTE CHE RISOLVE IL TUO PROBLEMA ---
	if err == sql.ErrNoRows {
		// Non abbiamo trovato nessuno con questo nome.
		// È OTTIMO se stiamo controllando la disponibilità.
		// Restituiamo id vuoto e nil (nessun errore).
		return "", nil
	}
	// -----------------------------------------------------

	if err != nil {
		// Questo è un errore vero (es. database spento)
		return "", err
	}

	return id, nil
}
