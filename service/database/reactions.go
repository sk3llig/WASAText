package database

import (
	"database/sql"
	"fmt"
)

// Definiscila qui, locale, così non litiga con gli altri pacchetti
type ReactionDetail struct {
	Emoji    string `json:"emoji"`
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}

// GetMessageReactionsList restituisce la lista dettagliata per il modale
func (db *appdbimpl) GetMessageReactionsList(msgID string) ([]ReactionDetail, error) {
	// Query mirata: JOIN tra reactions e users
	query := `
        SELECT r.emoji, u.id, u.name,
        FROM reactions r
        JOIN users u ON r.user_id = u.id
        WHERE r.message_id = ?
    `

	rows, err := db.c.Query(query, msgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ReactionDetail

	for rows.Next() {
		var r ReactionDetail
		var photoSql sql.NullString // Gestiamo il caso in cui la foto sia null

		// Scansioniamo i dati
		if err := rows.Scan(&r.Emoji, &r.UserID, &r.UserName, &photoSql); err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	return results, nil
}

// AddReaction aggiunge una reazione (emoji) a un messaggio
func (db *appdbimpl) AddReaction(convID, msgID, userID, emoji string) error {
	// PRIMA ERA: INSERT INTO reactions ...
	// ORA DIVENTA: INSERT OR IGNORE INTO reactions ...

	query := `INSERT OR IGNORE INTO reactions (conversation_id, message_id, user_id, emoji) VALUES (?, ?, ?, ?)`

	_, err := db.c.Exec(query, convID, msgID, userID, emoji)
	if err != nil {
		return fmt.Errorf("error adding reaction: %w", err)
	}
	return nil
}

// RemoveReaction rimuove una reazione specifica
func (db *appdbimpl) RemoveReaction(convID, msgID, userID, emoji string) error {
	var query string
	var err error

	if emoji != "" {
		// Rimuove solo quell'emoji specifica
		query = `DELETE FROM reactions WHERE message_id = ? AND user_id = ? AND emoji = ?`
		_, err = db.c.Exec(query, msgID, userID, emoji)
	} else {
		// (Opzionale) Se l'emoji è vuota, potremmo decidere di rimuovere TUTTE le reazioni di quell'utente
		// Ma per sicurezza implementiamo la rimozione puntuale come da specifica
		query = `DELETE FROM reactions WHERE message_id = ? AND user_id = ?`
		_, err = db.c.Exec(query, msgID, userID)
	}

	if err != nil {
		return fmt.Errorf("error removing reaction: %w", err)
	}
	return nil
}
