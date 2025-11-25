package database

import (
	"errors"
	"fmt"
)

// AddToGroup aggiunge un utente alla tabella participants
func (db *appdbimpl) AddToGroup(groupID string, userID string) error {
	query := `INSERT INTO participants (conversation_id, user_id) VALUES (?, ?)`
	_, err := db.c.Exec(query, groupID, userID)
	if err != nil {
		return fmt.Errorf("error adding participant: %w", err)
	}
	return nil
}

// LeaveGroup rimuove un utente dalla tabella participants
func (db *appdbimpl) LeaveGroup(groupID string, userID string) error {
	query := `DELETE FROM participants WHERE conversation_id = ? AND user_id = ?`
	_, err := db.c.Exec(query, groupID, userID)
	if err != nil {
		return fmt.Errorf("error removing participant: %w", err)
	}
	return nil
}

// SetGroupName aggiorna il nome (solo se è un gruppo)
func (db *appdbimpl) SetGroupName(groupID string, name string) error {
	// Opzionale: aggiungere AND is_group=1 per sicurezza
	query := `UPDATE conversations SET name = ? WHERE id = ? AND is_group = 1`
	_, err := db.c.Exec(query, name, groupID)
	return err
}

// SetGroupPhoto aggiorna la foto (solo se è un gruppo)
func (db *appdbimpl) SetGroupPhoto(groupID string, photoURL string) error {
	query := `UPDATE conversations SET photo_url = ? WHERE id = ? AND is_group = 1`
	_, err := db.c.Exec(query, photoURL, groupID)
	return err
}

// GetGroupMembers restituisce gli ID degli utenti che fanno parte di un gruppo
func (db *appdbimpl) GetGroupMembers(groupID string) ([]string, error) {
	query := `SELECT user_id FROM participants WHERE conversation_id = ?`
	rows, err := db.c.Query(query, groupID)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	var members []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		members = append(members, uid)
	}
	return members, nil
}

// DeleteGroup elimina fisicamente il gruppo (cascade cancellerà i messaggi e i partecipanti)
func (db *appdbimpl) DeleteGroup(groupID string) error {
	// Controllo opzionale: verificare che sia un gruppo
	query := `DELETE FROM conversations WHERE id = ? AND is_group = 1`
	res, err := db.c.Exec(query, groupID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("group not found")
	}
	return nil
}
