package api

import "time"

// -------------------------------------------------------
// RISPOSTE STANDARD
// -------------------------------------------------------

// SuccessMessage rappresenta una risposta generica di successo (200/204)
type SuccessMessage struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ErrorDetail rappresenta il dettaglio di un errore in un campo specifico
type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error rappresenta il formato standard di errore (4xx/5xx)
type Error struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Errors  []ErrorDetail `json:"errors,omitempty"` // Omette il campo se vuoto
	Trace   string        `json:"trace,omitempty"`  // ID per il debug
}

// -------------------------------------------------------
// AUTHENTICATION
// -------------------------------------------------------

type LoginRequest struct {
	Name string `json:"name"`
}

type LoginResponse struct {
	Identifier string `json:"identifier"`
}

// -------------------------------------------------------
// MEDIA (Upload Immagini)
// -------------------------------------------------------

// ImageMetadata è la risposta dopo un upload su /images
type ImageMetadata struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
}

// -------------------------------------------------------
// UTENTI & PROFILO
// -------------------------------------------------------

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	PhotoURL string `json:"photo_url"`
}

type SetUserNameRequest struct {
	Name string `json:"name"`
}

// SetPhotoRequest accetta un URL (stringa), non più il file binario
type SetPhotoRequest struct {
	PhotoURL string `json:"photo_url"`
}

// -------------------------------------------------------
// CONVERSAZIONI
// -------------------------------------------------------

type Conversation struct {
	ID              string    `json:"id"`
	Name            *string   `json:"name"` // Puntatore perché può essere null (JSON null)
	IsGroup         bool      `json:"is_group"`
	PhotoURL        string    `json:"photo_url"`
	LastMessage     string    `json:"last_message"`
	LastMessageTime time.Time `json:"last_message_time"`
}

type ConversationsList struct {
	Conversations []Conversation `json:"conversations"`
}

type CreateConversationRequest struct {
	ParticipantIDs []string `json:"participant_ids"`
	IsGroup        bool     `json:"is_group"`
	GroupName      string   `json:"group_name,omitempty"`
}

type ConversationDetails struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	IsGroup      bool      `json:"is_group"`
	Participants []User    `json:"participants"`
	Messages     []Message `json:"messages"`
}

// -------------------------------------------------------
// MESSAGGI
// -------------------------------------------------------

type Message struct {
	ID             string     `json:"id"`
	ConversationID string     `json:"conversation_id"`
	Sender         User       `json:"sender"`
	Timestamp      time.Time  `json:"timestamp"`
	Content        string     `json:"content"`
	ReplyTo        string     `json:"reply_to,omitempty"`
	Type           string     `json:"type"`        // "text" o "photo"
	Attachments    []string   `json:"attachments"` // Array di URL
	Status         string     `json:"status"`      // "sent", "received", "read
	Reactions      []Reaction `json:"reactions"`
}

type MessagesList struct {
	Messages []Message `json:"messages"`
}

type SendMessageRequest struct {
	Content     string   `json:"content,omitempty"`
	ReplyTo     string   `json:"reply_to,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}

type ForwardMessageRequest struct {
	TargetConversationID string `json:"target_conversation_id"`
	Comment              string `json:"comment,omitempty"`
}

// -------------------------------------------------------
// REAZIONI
// -------------------------------------------------------

type CommentRequest struct {
	Emoji string `json:"emoji"`
}

type Reaction struct {
	ID    string `json:"id"`
	Emoji string `json:"emoji"`
	User  User   `json:"user"`
}

// -------------------------------------------------------
// GRUPPI
// -------------------------------------------------------

type Group struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Members   []string  `json:"members"` // Lista di ID utente
	PhotoURL  string    `json:"photo_url"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupsList struct {
	Groups []Group `json:"groups"`
}

type CreateGroupRequest struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

type AddToGroupRequest struct {
	UserID string `json:"user_id"`
}

type SetGroupNameRequest struct {
	Name string `json:"name"`
}

type MembersList struct {
	Members []User `json:"members"`
}
