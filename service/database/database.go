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
	ID       string
	Name     string // Vuoto per chat 1-to-1
	IsGroup  bool
	PhotoURL string

	CreatorID string // Essenziale per sapere chi può rimuovere i membri
}

type Message struct {
	ID             string
	ConversationID string
	SenderID       string
	SenderName     string // Popolato via JOIN
	Content        string
	Timestamp      time.Time
	Type           string
	Attachments    []string
	ReplyTo        string
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
		FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
		FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE SET NULL
	);`)
	if err != nil {
		return nil, fmt.Errorf("error creating messages table: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
