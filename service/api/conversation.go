package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// createConversation (POST /conversations)
func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 1. Auth
	userID, err := rt.extractUser(r)
	if err != nil {
		rt.sendJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// 2. Parse Body
	var req CreateConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.sendJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// 3. Validazione
	// Aggiungiamo chi crea la chat alla lista partecipanti (se non c'è già)
	hasMe := false
	for _, p := range req.ParticipantIDs {
		if p == userID {
			hasMe = true
			break
		}
	}
	if !hasMe {
		req.ParticipantIDs = append(req.ParticipantIDs, userID)
	}

	if len(req.ParticipantIDs) < 2 {
		rt.sendJSONError(w, http.StatusBadRequest, "At least 2 participants required")
		return
	}

	// TODO: Qui dovresti controllare se è una chat 1-to-1 e se esiste già.
	// Per l'esame base, possiamo saltare questo controllo complesso o implementarlo dopo.

	// 4. Crea Oggetto DB
	newConv := database.Conversation{
		ID:      uuid.New().String(),
		IsGroup: req.IsGroup,
	}
	if req.IsGroup {
		newConv.Name = req.GroupName
	}

	// 5. Salva nel DB
	savedConv, err := rt.db.CreateConversation(newConv, req.ParticipantIDs)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error creating conversation")
		rt.sendJSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// 6. Risposta API
	apiConv := Conversation{
		ID:      savedConv.ID,
		IsGroup: savedConv.IsGroup,
		Name:    &savedConv.Name, // Puntatore per gestire null
	}
	// Se è vuota (1-to-1), Name deve essere nil nel JSON
	if savedConv.Name == "" {
		apiConv.Name = nil
	}

	rt.sendJSONResponse(w, http.StatusCreated, apiConv)
}

// getMyConversations (GET /me/conversations)
func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID, err := rt.extractUser(r)
	if err != nil {
		rt.sendJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Recupera dal DB
	dbConvs, err := rt.db.GetConversationsByUser(userID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error fetching conversations")
		rt.sendJSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Converti
	var res ConversationsList
	res.Conversations = make([]Conversation, 0)

	for _, c := range dbConvs {
		apiC := Conversation{
			ID:       c.ID,
			IsGroup:  c.IsGroup,
			PhotoURL: c.PhotoURL,
		}
		if c.Name != "" {
			nameCopy := c.Name
			apiC.Name = &nameCopy
		}
		res.Conversations = append(res.Conversations, apiC)
	}

	rt.sendJSONResponse(w, http.StatusOK, res)
}
