package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

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

	// Controlliamo se TUTTI gli utenti invitati esistono davvero
	// ================== BLOCCO INTELLIGENTE ==================
	// Ciclo su tutti i partecipanti richiesti
	for i, pID := range req.ParticipantIDs {
		// 1. Proviamo a vedere se è un ID valido
		exists, _ := rt.db.CheckUserExists(pID)

		if exists {
			// È un ID valido, tutto ok! Continua col prossimo.
			continue
		}

		// 2. Se non è un ID, forse è un NOME? Proviamo a cercare l'ID dal nome.
		userObj, err := rt.db.GetUserByName(pID)
		if err != nil {
			// Non è né un ID né un Nome valido -> ERRORE
			//rt.baseLogger.WithError(err).Error("User not found by name or ID")
			rt.sendJSONError(w, http.StatusBadRequest, "User not found: "+pID)
			return
		}

		// 3. Trovato! Sostituiamo il nome ("Alice") con l'ID vero ("5af5...")
		// così il resto del codice e il database sono felici.
		req.ParticipantIDs[i] = userObj.ID
	}
	// =========================================================
	// ================== [FINE NUOVO BLOCCO] ====================

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
		Name:    &savedConv.Name,
	}
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

// getMyGroups (GET /me/groups)
func (rt *_router) getMyGroups(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID, err := rt.extractUser(r)
	if err != nil {
		rt.sendJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Riutilizziamo la query che prende tutte le chat e filtriamo in Go
	allConvs, err := rt.db.GetConversationsByUser(userID)
	if err != nil {
		rt.sendJSONError(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	var res GroupsList
	res.Groups = make([]Group, 0)

	for _, c := range allConvs {
		if c.IsGroup {
			res.Groups = append(res.Groups, Group{
				ID:       c.ID,
				Name:     c.Name,
				PhotoURL: c.PhotoURL,
				// Members: []string{} // Non carichiamo i membri nella lista per velocità
			})
		}
	}
	rt.sendJSONResponse(w, http.StatusOK, res)
}

// getConversation (GET /conversations/:conversation_id)
// deve tornare Info + Partecipanti + Ultimi Messaggi
func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_, err := rt.extractUser(r)
	if err != nil {
		rt.sendJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	convID := ps.ByName("conversation_id")

	// 1. Info Conversazione
	dbConv, err := rt.db.GetConversationByID(convID)
	if err != nil {
		rt.sendJSONError(w, http.StatusNotFound, "Conversation not found")
		return
	}

	// 2. Lista Messaggi
	dbMsgs, err := rt.db.GetMessagesByConversation(convID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error fetching messages") // Non bloccante, magari lista vuota
	}

	// 3. Recupero Partecipanti (Serve una query GetGroupMembers o simile)
	// Nota: Se è 1-to-1 usiamo GetGroupMembers lo stesso, funziona su 'participants'
	memberIDs, err := rt.db.GetGroupMembers(convID)
	if err != nil {
		// Logghiamo l'errore ma non blocchiamo tutto (magari torniamo la chat senza membri)
		// Oppure ritorniamo errore 500. Scegliamo la via sicura (500)
		rt.baseLogger.WithError(err).Error("Error fetching group members")
		rt.sendJSONError(w, http.StatusInternalServerError, "Internal Error")
		return
	}
	// Costruzione risposta complessa
	details := ConversationDetails{
		ID:           dbConv.ID,
		Name:         dbConv.Name,
		IsGroup:      dbConv.IsGroup,
		Participants: []User{},
		Messages:     []Message{},
	}

	// Popola partecipanti (con dettagli User)
	for _, mid := range memberIDs {
		u, _ := rt.db.GetUserByID(mid)
		details.Participants = append(details.Participants, User{ID: u.ID, Name: u.Name, PhotoURL: u.PhotoURL})
	}

	// Popola messaggi
	for _, m := range dbMsgs {
		details.Messages = append(details.Messages, Message{
			ID: m.ID, ConversationID: m.ConversationID, Content: m.Content,
			Sender:    User{ID: m.SenderID, Name: m.SenderName},
			Timestamp: m.Timestamp, Type: m.Type, Attachments: m.Attachments,
		})
	}

	rt.sendJSONResponse(w, http.StatusOK, details)
}

// deleteConversation (DELETE /conversations/:id)
func (rt *_router) deleteConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_, err := rt.extractUser(r)
	if err != nil {
		rt.sendJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	// TODO: Controllare permessi (solo admin/owner o partecipante?)

	err = rt.db.DeleteConversation(ps.ByName("conversation_id"))
	if err != nil {
		rt.sendJSONError(w, http.StatusNotFound, "Conversation not found")
		return
	}
	rt.sendJSONResponse(w, http.StatusOK, SuccessMessage{Success: true, Message: "Deleted"})
}
