package database

import "fmt"

// AddReaction aggiunge una reazione (emoji) a un messaggio
func (db *appdbimpl) AddReaction(convID, msgID, userID, emoji string) error {
	// Vincolo di unicità: la Primary Key (message_id, user_id, emoji) impedisce duplicati.
	// Se l'utente prova a mettere due volte la stessa emoji, SQLite darà errore,
	// che noi restituiamo al chiamante.
	query := `INSERT INTO reactions (conversation_id, message_id, user_id, emoji) VALUES (?, ?, ?, ?)`

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
