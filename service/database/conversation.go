package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// CreateConversation crea la chat e associa i partecipanti
func (db *appdbimpl) CreateConversation(c Conversation, participantIDs []string) (Conversation, error) {
	// 1. Iniziamo una Transazione (Atomic Operation)
	tx, err := db.c.Begin()
	if err != nil {
		return Conversation{}, err
	}
	// Se qualcosa va male, facciamo Rollback alla fine
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 2. Inseriamo la conversazione
	queryConv := `INSERT INTO conversations (id, name, is_group, photo_url) VALUES (?, ?, ?, ?)`
	if _, err = tx.Exec(queryConv, c.ID, c.Name, c.IsGroup, c.PhotoURL); err != nil {
		return Conversation{}, fmt.Errorf("error inserting conversation: %w", err)
	}

	// 3. Inseriamo i partecipanti uno ad uno
	queryPart := `INSERT INTO participants (conversation_id, user_id) VALUES (?, ?)`
	for _, userID := range participantIDs {
		if _, err = tx.Exec(queryPart, c.ID, userID); err != nil {
			return Conversation{}, fmt.Errorf("error adding participant %s: %w", userID, err)
		}
	}

	// 4. Confermiamo tutto (Commit)
	if err = tx.Commit(); err != nil {
		return Conversation{}, err
	}

	return c, nil
}

// GetConversationsByUser trova tutte le chat dove l'utente è presente
// GetConversationsByUser ritorna le chat dell'utente con la logica "Smart Name/Photo"

func (db *appdbimpl) GetConversationsByUser(userID string) ([]Conversation, error) {
	// Questa query fa la "magia":
	// 1. Seleziona le conversazioni (c) dove io (p_me) sono partecipante.
	// 2. Cerca un "altro" partecipante (p_other) che NON sono io.
	// 3. Unisce con la tabella utenti (u) per avere i dati dell'altro.
	// 4. CASE WHEN: Se è un gruppo usa i dati di 'c', altrimenti usa i dati di 'u'.

	query := `
		SELECT
			c.id,
			c.is_group,
			CASE
				WHEN c.is_group = 1 THEN c.name
				ELSE u.name
			END AS display_name,
			CASE
				WHEN c.is_group = 1 THEN c.photo_url
				ELSE u.photo_url
			END AS display_photo
		FROM conversations c
		JOIN participants p_me ON c.id = p_me.conversation_id
		LEFT JOIN participants p_other ON c.id = p_other.conversation_id AND p_other.user_id != ?
		LEFT JOIN users u ON p_other.user_id = u.id
		WHERE p_me.user_id = ?
		GROUP BY c.id
	`

	// Passiamo userID due volte (una per escluderlo dal p_other, una per il WHERE p_me)
	rows, err := db.c.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	defer rows.Close()

	var res []Conversation
	for rows.Next() {
		var c Conversation
		var name sql.NullString
		var photo sql.NullString

		// Leggiamo i dati calcolati dall'SQL
		if err := rows.Scan(&c.ID, &c.IsGroup, &name, &photo); err != nil {
			return nil, err
		}

		if name.Valid {
			c.Name = name.String
		}
		if photo.Valid {
			c.PhotoURL = photo.String
		}
		res = append(res, c)
	}
	return res, nil
}

// GetConversationByID recupera i dettagli di una chat specifica
func (db *appdbimpl) GetConversationByID(convID string) (Conversation, error) {
	var c Conversation
	query := `SELECT id, name, is_group, photo_url FROM conversations WHERE id = ?`
	err := db.c.QueryRow(query, convID).Scan(&c.ID, &c.Name, &c.IsGroup, &c.PhotoURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Conversation{}, errors.New("conversation not found")
		}
		return Conversation{}, err
	}
	return c, nil
}

func (db *appdbimpl) DeleteConversation(convID string) error {
	query := `DELETE FROM conversations WHERE id = ?`
	res, err := db.c.Exec(query, convID)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return errors.New("conversation not found")
	}
	return nil
}
