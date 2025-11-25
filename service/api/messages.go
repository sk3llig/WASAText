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

// getMessage (GET .../messages/:message_id)
func (rt *_router) getMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_, err := rt.extractUser(r) // Auth check
	if err != nil {
		rt.sendJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	m, err := rt.db.GetMessageByID(ps.ByName("message_id"))
	if err != nil {
		rt.sendJSONError(w, http.StatusNotFound, "Message not found")
		return
	}

	// Converti DB -> API
	apiMsg := Message{
		ID: m.ID, ConversationID: m.ConversationID, Content: m.Content,
		Sender:    User{ID: m.SenderID, Name: m.SenderName},
		Timestamp: m.Timestamp, Type: m.Type, Attachments: m.Attachments, ReplyTo: m.ReplyTo,
	}
	rt.sendJSONResponse(w, http.StatusOK, apiMsg)
}

// deleteMessage (DELETE .../messages/:message_id)
func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, err := rt.extractUser(r)
	if err != nil {
		rt.sendJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// DeleteMessage nel DB deve controllare che userID sia il sender!
	err = rt.db.DeleteMessage(ps.ByName("conversation_id"), ps.ByName("message_id"), userID)
	if err != nil {
		// Se l'errore è "non trovato o non autore", torniamo 403 o 404
		rt.sendJSONError(w, http.StatusForbidden, "Cannot delete: not found or not author")
		return
	}
	rt.sendJSONResponse(w, http.StatusOK, SuccessMessage{Success: true, Message: "Deleted"})
}

// forwardMessage (POST .../forward)
func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, err := rt.extractUser(r)
	if err != nil {
		rt.sendJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// 1. Leggi il messaggio originale
	originalMsgID := ps.ByName("message_id")
	originalMsg, err := rt.db.GetMessageByID(originalMsgID)
	if err != nil {
		rt.sendJSONError(w, http.StatusNotFound, "Original message not found")
		return
	}

	// 2. Leggi dove mandarlo
	var req ForwardMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.sendJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// 3. Crea il nuovo messaggio (copia contenuto/allegati)
	// Recuperiamo info utente per la risposta rapida
	me, _ := rt.db.GetUserByID(userID)

	newMsg := database.Message{
		ID:             uuid.New().String(),
		ConversationID: req.TargetConversationID, // Nuova destinazione
		SenderID:       userID,                   // Io sono il nuovo mittente
		SenderName:     me.Name,
		Content:        originalMsg.Content,     // Copia contenuto
		Type:           originalMsg.Type,        // Copia tipo
		Attachments:    originalMsg.Attachments, // Copia allegati
		Timestamp:      time.Now(),
	}

	// 4. Salva
	_, err = rt.db.CreateMessage(newMsg)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Forward failed")
		rt.sendJSONError(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	// Rispondi col nuovo messaggio
	apiMsg := Message{
		ID: newMsg.ID, ConversationID: newMsg.ConversationID, Sender: User{ID: userID, Name: me.Name},
		Content: newMsg.Content, Timestamp: newMsg.Timestamp, Type: newMsg.Type, Attachments: newMsg.Attachments,
		Status: "sent",
	}
	rt.sendJSONResponse(w, http.StatusCreated, apiMsg)
}

// commentMessage (POST .../reactions)
func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, err := rt.extractUser(r)
	if err != nil {
		rt.sendJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.sendJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// AddReaction deve essere implementato nel DB (Tabella reactions)
	err = rt.db.AddReaction(ps.ByName("conversation_id"), ps.ByName("message_id"), userID, req.Emoji)
	if err != nil {
		rt.sendJSONError(w, http.StatusInternalServerError, "Error adding reaction")
		return
	}

	// Mock response (in realtà dovresti ricaricare la reazione)
	user, _ := rt.db.GetUserByID(userID)
	rt.sendJSONResponse(w, http.StatusCreated, Reaction{Emoji: req.Emoji, User: User{ID: user.ID, Name: user.Name}})
}

// uncommentMessage (DELETE .../reactions)
func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, err := rt.extractUser(r)
	if err != nil {
		rt.sendJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	emoji := r.URL.Query().Get("emoji")
	err = rt.db.RemoveReaction(ps.ByName("conversation_id"), ps.ByName("message_id"), userID, emoji)
	if err != nil {
		rt.sendJSONError(w, http.StatusInternalServerError, "Error removing reaction")
		return
	}
	rt.sendJSONResponse(w, http.StatusOK, SuccessMessage{Success: true, Message: "Reaction removed"})
}
