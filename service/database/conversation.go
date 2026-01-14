package database

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
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
			_ = tx.Rollback()
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
	query := `
SELECT
    c.id,
    c.is_group,
    CASE WHEN c.is_group = 1 THEN c.name ELSE u.name END AS display_name,
    CASE WHEN c.is_group = 1 THEN c.photo_url ELSE u.photo_url END AS display_photo,
    COALESCE((SELECT content FROM messages WHERE conversation_id = c.id ORDER BY timestamp DESC LIMIT 1), '') as last_msg,
    (SELECT timestamp FROM messages WHERE conversation_id = c.id ORDER BY timestamp DESC LIMIT 1) as last_time,
    (SELECT COUNT(*) FROM messages WHERE conversation_id = c.id AND sender_id != ? AND status != 'read') as unread_cnt
FROM conversations c
JOIN participants p_me ON c.id = p_me.conversation_id
LEFT JOIN participants p_other ON c.id = p_other.conversation_id AND p_other.user_id != ?
LEFT JOIN users u ON p_other.user_id = u.id
WHERE p_me.user_id = ?
GROUP BY c.id`

	rows, err := db.c.Query(query, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var res []Conversation
	for rows.Next() {
		var c Conversation
		var name, photo, lastMsg sql.NullString
		var lastTime sql.NullTime

		if err := rows.Scan(&c.ID, &c.IsGroup, &name, &photo, &lastMsg, &lastTime, &c.UnreadCount); err != nil {
			return nil, err
		}

		c.Name = name.String
		c.PhotoURL = photo.String
		c.LastMessage = lastMsg.String
		if lastTime.Valid {
			c.LastMessageTime = lastTime.Time
		}

		// Gestione icone per messaggi vuoti (foto)
		if lastMsg.Valid && lastMsg.String == "" {
			c.LastMessage = "📷 Foto"
		}

		res = append(res, c)
	}
	// Inserisci questo prima di 'return res, nil' nella funzione del DB
	sort.Slice(res, func(i, j int) bool {
		// Se una chat non ha messaggi, la mettiamo in fondo
		if res[i].LastMessageTime.IsZero() {
			return false
		}
		if res[j].LastMessageTime.IsZero() {
			return true
		}
		// Ordine decrescente: la più recente (più grande) prima
		return res[i].LastMessageTime.After(res[j].LastMessageTime)
	})
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
