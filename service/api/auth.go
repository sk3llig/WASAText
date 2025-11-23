package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// doLogin gestisce il login (POST /session).
// Se l'utente esiste, lo logga. Se non esiste, lo crea.
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 1. Decodifica il body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.sendJSONError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// 2. Validazione minima (come da YAML: min 3, max 32)
	if len(req.Name) < 3 || len(req.Name) > 32 {
		rt.sendJSONError(w, http.StatusBadRequest, "Username must be between 3 and 32 chars")
		return
	}

	// 3. Cerchiamo l'utente nel DB
	user, err := rt.db.GetUserByName(req.Name)
	if err != nil {
		// Se l'errore NON è "non trovato", allora è un errore grave del DB
		if err.Error() != "user not found" {
			rt.baseLogger.WithError(err).Error("Database error checking user")
			rt.sendJSONError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		// 4. Utente non trovato -> CREIAMO UN NUOVO UTENTE
		newUser := database.User{
			ID:   uuid.New().String(), // Generiamo un nuovo UUID
			Name: req.Name,
		}

		createdUser, errCreate := rt.db.CreateUser(newUser)
		if errCreate != nil {
			rt.baseLogger.WithError(errCreate).Error("Database error creating user")
			rt.sendJSONError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		// Assegniamo l'utente creato alla variabile user per la risposta
		user = createdUser

		// Loggiamo l'evento
		rt.baseLogger.Infof("Created new user: %s (%s)", user.Name, user.ID)
	}

	// 5. Risposta (201 Created)
	// Ritorniamo l'identificativo (Token)
	resp := LoginResponse{
		Identifier: user.ID,
	}

	rt.sendJSONResponse(w, http.StatusCreated, resp)
}
