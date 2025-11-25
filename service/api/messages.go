package api

import (
	"encoding/json"
	"net/http"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// sendMessage (POST /conversations/:id/messages)
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, err := rt.extractUser(r)
	if err != nil {
		rt.sendJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	convID := ps.ByName("conversation_id")

	// Parse input
	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.sendJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Validazione: O testo O allegati
	if req.Content == "" && len(req.Attachments) == 0 {
		rt.sendJSONError(w, http.StatusBadRequest, "Message must have content or attachments")
		return
	}

	// Tipo messaggio
	msgType := "text"
	if len(req.Attachments) > 0 {
		msgType = "photo"
	}

	// Recupera nome sender per la risposta
	user, _ := rt.db.GetUserByID(userID)

	// Prepara oggetto DB
	dbMsg := database.Message{
		ID:             uuid.New().String(),
		ConversationID: convID,
		SenderID:       userID,
		SenderName:     user.Name, // Usato solo per risposta immediata
		Content:        req.Content,
		Timestamp:      time.Now(),
		Type:           msgType,
		Attachments:    req.Attachments,
		ReplyTo:        req.ReplyTo,
	}

	// Salva
	_, err = rt.db.CreateMessage(dbMsg)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error saving message")
		rt.sendJSONError(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	// Rispondi (Mappiamo DB struct -> API struct)
	apiMsg := Message{
		ID:             dbMsg.ID,
		ConversationID: dbMsg.ConversationID,
		Sender: User{
			ID:   dbMsg.SenderID,
			Name: user.Name, // Nome fresco dal DB
		},
		Timestamp:   dbMsg.Timestamp,
		Content:     dbMsg.Content,
		Type:        dbMsg.Type,
		Attachments: dbMsg.Attachments,
		ReplyTo:     dbMsg.ReplyTo,
		Status:      "sent",
	}

	rt.sendJSONResponse(w, http.StatusCreated, apiMsg)
}

// getConversationMessages (GET /conversations/:id/messages)
func (rt *_router) getConversationMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_, err := rt.extractUser(r)
	if err != nil {
		rt.sendJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	convID := ps.ByName("conversation_id")

	// TODO: Controllare se userID fa parte della conversation (Sicurezza)
	// es: if !rt.db.IsUserInConversation(convID, userID) { Error 403 }

	msgs, err := rt.db.GetMessagesByConversation(convID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting messages")
		rt.sendJSONError(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	// Mappatura DB -> API
	var res MessagesList
	res.Messages = make([]Message, 0)
	for _, m := range msgs {
		res.Messages = append(res.Messages, Message{
			ID:             m.ID,
			ConversationID: m.ConversationID,
			Sender: User{
				ID:   m.SenderID,
				Name: m.SenderName, // Arriva dalla JOIN
			},
			Timestamp:   m.Timestamp,
			Content:     m.Content,
			Type:        m.Type,
			Attachments: m.Attachments,
			ReplyTo:     m.ReplyTo,
			Status:      "read",
		})
	}

	rt.sendJSONResponse(w, http.StatusOK, res)
}
