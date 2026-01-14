package database

import (
	"database/sql"
	"fmt"
	"time"
)

// ---------------------------------------------------------
// 1. STRUTTURE DATI INTERNE (Lato DB)
// ---------------------------------------------------------

type User struct {
	ID       string
	Name     string
	PhotoURL string
}

type Conversation struct {
	// Campi che esistono nella tabella del DB
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsGroup  bool   `json:"is_group"`
	PhotoURL string `json:"photo_url"`

	// CreatorID serve internamente, ma spesso non si manda al frontend nella lista.
	// Se ti serve, tienilo, ma aggiungi il tag json se vuoi che il frontend lo veda.
	CreatorID string `json:"creator_id,omitempty"`

	// --- CAMPI "CALCOLATI" (Non esistono nella tabella conversations, ma servono alla Struct) ---
	// Go deve sapere che questi campi esistono per poterci mettere dentro
	// i risultati delle query che abbiamo scritto prima.
	LastMessage     string    `json:"last_message"`
	LastMessageTime time.Time `json:"last_message_time"`
	UnreadCount     int       `json:"unread_count"`
}

// NUOVA STRUCT REACTION
type ReactionSimple struct {
	Emoji    string `json:"emoji"`
	UserName string `json:"user_name"` // Solo il nome, come volevi
}

type Message struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	SenderID       string    `json:"sender_id"`
	SenderName     string    `json:"sender_name"` // IL TUO VECCHIO CAMPO (INVARIATO)
	Content        string    `json:"content"`
	Timestamp      time.Time `json:"timestamp"`
	Type           string    `json:"type"`
	Attachments    []string  `json:"attachments"`
	ReplyTo        string    `json:"reply_to,omitempty"`
	Status         string
	// UNICA AGGIUNTA ALLA STRUCT:
	Reactions []ReactionSimple `json:"reactions"`
}

// ---------------------------------------------------------
// 2. INTERFACCIA
// ---------------------------------------------------------

type AppDatabase interface {
	Ping() error

	// --- USERS ---
	CreateUser(u User) (User, error)
	GetUserByName(name string) (User, error)
	GetUserByID(id string) (User, error)
	SetUserName(id string, name string) error
	SetPhotoUrl(id string, url string) error
	CheckUserByUsername(username string) (bool, error)
	GetUserIDByName(name string) (string, error)

	// --- CONVERSATIONS ---
	CreateConversation(conv Conversation, participantIDs []string) (Conversation, error)
	GetConversationsByUser(userID string) ([]Conversation, error)
	GetConversationByID(id string) (Conversation, error)
	DeleteConversation(convID string) error
	CheckUserExists(id string) (bool, error)

	// --- MESSAGES ---
	CreateMessage(m Message) (Message, error)
	GetMessagesByConversation(convID string) ([]Message, error)
	GetMessageByID(msgID string) (Message, error)
	DeleteMessage(convID, msgID, userID string) error
	MarkAsRead(conversationID string, userID string) error

	// --- GROUPS ---
	AddToGroup(groupID string, userID string) error
	LeaveGroup(groupID string, userID string) error
	SetGroupName(groupID string, name string) error
	SetGroupPhoto(groupID string, photoURL string) error
	GetGroupMembers(groupID string) ([]string, error)
	DeleteGroup(groupID string) error

	// --- REACTIONS (EXTRA) ---
	AddReaction(convID, msgID, userID, emoji string) error
	RemoveReaction(convID, msgID, userID, emoji string) error
}

// ---------------------------------------------------------
// 3. IMPLEMENTAZIONE
// ---------------------------------------------------------

type appdbimpl struct {
	c *sql.DB
}

// New crea una nuova istanza del database
func New(db *sql.DB) (*appdbimpl, error) {
	if err := db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, fmt.Errorf("error enabling foreign keys: %w", err)
	}

	// --- CREAZIONE TABELLE ---

	// 1. Users
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id TEXT NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		photo_url TEXT
	);`)
	if err != nil {
		return nil, fmt.Errorf("error creating users table: %w", err)
	}

	// 2. Conversations
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS conversations (
		id TEXT NOT NULL PRIMARY KEY,
		name TEXT,
		is_group BOOLEAN NOT NULL,
		photo_url TEXT
	);`)
	if err != nil {
		return nil, fmt.Errorf("error creating conversations table: %w", err)
	}

	// 3. Participants
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS participants (
		conversation_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		PRIMARY KEY (conversation_id, user_id),
		FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`)
	if err != nil {
		return nil, fmt.Errorf("error creating participants table: %w", err)
	}

	// 4. Messages
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id TEXT NOT NULL PRIMARY KEY,
		conversation_id TEXT NOT NULL,
		sender_id TEXT NOT NULL,
		content TEXT,
		timestamp DATETIME NOT NULL,
		type TEXT NOT NULL,
		attachments TEXT,
		reply_to TEXT,
		status TEXT NOT NULL DEFAULT 'received',
		FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
		FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE SET NULL
	);`)
	if err != nil {
		return nil, fmt.Errorf("error creating messages table: %w", err)
	}
	// 5. REACTIONS (QUESTA MANCAVA!)
	// Chiave primaria composta (message_id, user_id, emoji) per evitare duplicati
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS reactions (
       conversation_id TEXT NOT NULL,
       message_id TEXT NOT NULL,
       user_id TEXT NOT NULL,
       emoji TEXT NOT NULL,
       PRIMARY KEY (message_id, user_id, emoji),
       FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
       FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
       FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );`)
	if err != nil {
		return nil, fmt.Errorf("error creating reactions table: %w", err)
	}
	// 5. REACTIONS
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS reactions (
       conversation_id TEXT NOT NULL,
       message_id TEXT NOT NULL,
       user_id TEXT NOT NULL,
       emoji TEXT NOT NULL,

       -- Vincolo di unicità: Un utente può mettere 1 sola volta la stessa emoji sullo stesso messaggio
       PRIMARY KEY (message_id, user_id, emoji),

       -- Chiavi Esterne
       FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
       FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
       FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );`)
	if err != nil {
		return nil, fmt.Errorf("error creating reactions table: %w", err)
	}
	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
