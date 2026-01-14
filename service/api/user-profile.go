package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// getMyProfile restituisce il profilo dell'utente loggato
func (rt *_router) getMyProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 1. Identifica l'utente
	userID, err := rt.extractUser(r)
	if err != nil {
		rt.sendJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// 2. Recupera dati dal DB
	user, err := rt.db.GetUserByID(userID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Database error getting user")
		rt.sendJSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// 3. Risponde
	rt.sendJSONResponse(w, http.StatusOK, user)
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	userID, err := rt.extractUser(r)
	if err != nil {
		rt.sendJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req SetUserNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.sendJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// 1. Controllo Unicità
	ownerID, err := rt.db.GetUserIDByName(req.Name)
	if err != nil {
		// SE VEDI QUESTO ERRORE NEL TERMINALE, IL PROBLEMA È NELLA QUERY SQL
		rt.sendJSONError(w, http.StatusInternalServerError, "Database error checking name: "+err.Error())
		return
	}

	if ownerID != "" && ownerID != userID {
		rt.sendJSONError(w, http.StatusConflict, "Username taken")
		return
	}

	// 2. Aggiornamento
	err = rt.db.SetUserName(userID, req.Name)
	if err != nil {
		rt.sendJSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// 3. Risposta
	user, _ := rt.db.GetUserByID(userID)
	rt.sendJSONResponse(w, http.StatusOK, user)
}

// setMyPhoto aggiorna la foto profilo usando un URL (Upload-First pattern)
func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID, err := rt.extractUser(r)
	if err != nil {
		rt.sendJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req SetPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.sendJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Validazione URL (molto basilare, controlla solo se non è vuoto)
	if req.PhotoURL == "" {
		rt.sendJSONError(w, http.StatusBadRequest, "Photo URL is required")
		return
	}

	// Aggiorna DB
	err = rt.db.SetPhotoUrl(userID, req.PhotoURL)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Database error updating photo")
		rt.sendJSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Ritorna il profilo aggiornato
	user, _ := rt.db.GetUserByID(userID)
	rt.sendJSONResponse(w, http.StatusOK, user)
}

func (rt *_router) checkUserExistence(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")

	exists, err := rt.db.CheckUserByUsername(username)

	if err != nil {
		// --- AGGIUNGI QUESTA RIGA PER IL DEBUG ---
		// ----------------------------------------

		// Mandiamo l'errore anche al frontend per vederlo nella console del browser
		http.Error(w, "Errore interno server: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
