package database

import (
	"encoding/json"
	"fmt"
)

// CreateMessage salva un messaggio
func (db *appdbimpl) CreateMessage(m Message) (Message, error) {
	// Convertiamo l'array di allegati in stringa JSON per SQLite
	attachmentsJSON, err := json.Marshal(m.Attachments)
	if err != nil {
		// Se fallisce, salviamo array vuoto
		attachmentsJSON = []byte("[]")
	}

	query := `
		INSERT INTO messages (id, conversation_id, sender_id, content, timestamp, type, attachments, reply_to)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = db.c.Exec(query, m.ID, m.ConversationID, m.SenderID, m.Content, m.Timestamp, m.Type, string(attachmentsJSON), m.ReplyTo)
	if err != nil {
		return Message{}, fmt.Errorf("error inserting message: %w", err)
	}

	return m, nil
}

// GetMessagesByConversation recupera la cronologia con JOIN per il nome utente
func (db *appdbimpl) GetMessagesByConversation(convID string) ([]Message, error) {
	// Recuperiamo messaggi + Nome del mittente (JOIN users)
	query := `
		SELECT
			m.id, m.conversation_id, m.sender_id, u.name,
			m.content, m.timestamp, m.type, m.attachments, m.reply_to
		FROM messages m
		LEFT JOIN users u ON m.sender_id = u.id
		WHERE m.conversation_id = ?
		ORDER BY m.timestamp DESC
	`

	rows, err := db.c.Query(query, convID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var res []Message
	for rows.Next() {
		var m Message
		var attachmentsString string
		var replyToSql interface{} // Può essere NULL

		err = rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.SenderName, &m.Content, &m.Timestamp, &m.Type, &attachmentsString, &replyToSql)
		if err != nil {
			return nil, err
		}

		// Convertiamo JSON string -> Array
		_ = json.Unmarshal([]byte(attachmentsString), &m.Attachments)

		// Gestione reply_to nullable
		if replyToSql != nil {
			m.ReplyTo = replyToSql.(string)
		}

		res = append(res, m)
	}
	return res, nil
}
