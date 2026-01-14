package database

import (
	"database/sql"
	"encoding/json"
	"errors"
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

// GetMessageByID recupera un singolo messaggio (serve per getMessage e forward)
func (db *appdbimpl) GetMessageByID(msgID string) (Message, error) {
	var m Message
	var attachmentsString string
	var replyToSql interface{}

	// Join con users per avere il nome del mittente
	query := `
		SELECT m.id, m.conversation_id, m.sender_id, u.name, m.content, m.timestamp, m.type, m.attachments, m.reply_to
		FROM messages m
		LEFT JOIN users u ON m.sender_id = u.id
		WHERE m.id = ?
	`
	err := db.c.QueryRow(query, msgID).Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.SenderName, &m.Content, &m.Timestamp, &m.Type, &attachmentsString, &replyToSql)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Message{}, errors.New("message not found")
		}
		return Message{}, err
	}

	_ = json.Unmarshal([]byte(attachmentsString), &m.Attachments)
	if replyToSql != nil {
		m.ReplyTo = replyToSql.(string)
	}

	return m, nil
}

// GetMessagesByConversation recupera la cronologia con JOIN per il nome utente
func (db *appdbimpl) GetMessagesByConversation(convID string) ([]Message, error) {
	// 1. Query principale (Messaggi)
	query := `
       SELECT
          m.id, m.conversation_id, m.sender_id, u.name,
          m.content, m.timestamp, m.type, m.attachments, m.reply_to
       FROM messages m
       LEFT JOIN users u ON m.sender_id = u.id
       WHERE m.conversation_id = ?
       ORDER BY m.timestamp ASC
    `
	// Nota: Ho messo ASC (dal più vecchio al più nuovo) che è standard per le chat.
	// Se preferisci il contrario, rimetti DESC.

	rows, err := db.c.Query(query, convID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Message

	for rows.Next() {
		var m Message
		var attachmentsString string
		var replyToSql interface{}

		// 2. Lettura dati base
		err = rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.SenderName, &m.Content, &m.Timestamp, &m.Type, &attachmentsString, &replyToSql)
		if err != nil {
			return nil, err
		}

		_ = json.Unmarshal([]byte(attachmentsString), &m.Attachments)
		if replyToSql != nil {
			m.ReplyTo = replyToSql.(string)
		}

		// --- 3. PEZZO MANCANTE: RECUPERO REAZIONI ---
		// Senza questo blocco, m.Reactions rimane vuoto!
		rQuery := `
           SELECT r.emoji, u.name
           FROM reactions r
           JOIN users u ON r.user_id = u.id
           WHERE r.message_id = ?`

		rRows, rErr := db.c.Query(rQuery, m.ID)

		// Inizializziamo come slice vuota, così non è mai 'null' nel JSON
		var currentReactions = []ReactionSimple{}

		if rErr == nil {
			for rRows.Next() {
				var r ReactionSimple
				// Leggiamo emoji e nome utente
				if err := rRows.Scan(&r.Emoji, &r.UserName); err == nil {
					currentReactions = append(currentReactions, r)
				}
			}
			rRows.Close()
		}

		// Assegniamo la lista trovata al messaggio
		m.Reactions = currentReactions
		// -------------------------------------------

		res = append(res, m)
	}
	return res, nil
}

// DeleteMessage elimina un messaggio dal database
// Nota: Richiediamo userID per garantire che solo l'autore possa cancellarlo
func (db *appdbimpl) DeleteMessage(convID, msgID, userID string) error {
	// La query controlla sia l'ID del messaggio che il sender_id
	query := `DELETE FROM messages WHERE id = ? AND sender_id = ?`

	res, err := db.c.Exec(query, msgID, userID)
	if err != nil {
		return err
	}

	// Controlliamo se è stato cancellato qualcosa
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		// Se 0 righe, o il messaggio non esiste o l'utente non è l'autore
		return errors.New("message not found or user not authorized")
	}

	return nil
}
