package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Handler restituisce l'istanza di httprouter configurata
func (rt *_router) Handler() http.Handler {

	// --- AUTH ---
	rt.router.POST("/v1/session", rt.doLogin)

	// --- MEDIA ---
	rt.router.POST("/v1/images", rt.uploadImage)
	// ServeFiles permette di vedere le immagini dal browser (es. localhost:3000/images/foto.jpg)
	rt.router.ServeFiles("/images/*filepath", http.Dir("./images"))

	// --- PROFILE ---
	rt.router.GET("/v1/me", rt.getMyProfile)
	rt.router.PUT("/v1/me/name", rt.setMyUserName)
	rt.router.PUT("/v1/me/photo", rt.setMyPhoto)
	rt.router.GET("/v1/me/conversations", rt.getMyConversations)
	rt.router.GET("/v1/me/groups", rt.getMyGroups)

	// --- CONVERSATIONS ---
	rt.router.POST("/v1/conversations", rt.createConversation)
	rt.router.GET("/v1/conversations/:conversation_id", rt.getConversation)
	rt.router.DELETE("/v1/conversations/:conversation_id", rt.deleteConversation)

	// --- MESSAGES ---
	rt.router.GET("/v1/conversations/:conversation_id/messages", rt.getConversationMessages)
	rt.router.POST("/v1/conversations/:conversation_id/messages", rt.sendMessage)
	rt.router.GET("/v1/conversations/:conversation_id/messages/:message_id", rt.getMessage)
	rt.router.DELETE("/v1/conversations/:conversation_id/messages/:message_id", rt.deleteMessage)

	// --- MESSAGE ACTIONS ---
	rt.router.POST("/v1/conversations/:conversation_id/messages/:message_id/forward", rt.forwardMessage)
	rt.router.POST("/v1/conversations/:conversation_id/messages/:message_id/reactions", rt.commentMessage)
	rt.router.DELETE("/v1/conversations/:conversation_id/messages/:message_id/reactions", rt.uncommentMessage)

	// --- GROUPS ---
	rt.router.POST("/v1/groups", rt.createGroup)
	rt.router.GET("/v1/groups", rt.listGroups)
	rt.router.GET("/v1/groups/:group_id", rt.getGroup)
	rt.router.DELETE("/v1/groups/:group_id", rt.deleteGroup)

	// --- GROUP MEMBERS & DETAILS ---
	rt.router.GET("/v1/groups/:group_id/members", rt.listGroupMembers)
	rt.router.POST("/v1/groups/:group_id/members", rt.addToGroup)
	rt.router.DELETE("/v1/groups/:group_id/members/me", rt.leaveGroup)
	rt.router.PUT("/v1/groups/:group_id/name", rt.setGroupName)
	rt.router.PUT("/v1/groups/:group_id/photo", rt.setGroupPhoto)

	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}

// ----------------------------------------------------------------------
// PLACEHOLDERS (STUBS)
// Definiamo queste funzioni vuote per far compilare il progetto.
// Nelle prossime fasi le sposteremo nei file specifici (es. auth.go, messages.go)
// riempiendole con la logica vera.
// ----------------------------------------------------------------------

/*
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}
*/

/*
func (rt *_router) uploadImage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}
*/

/*
func (rt *_router) getMyProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}



func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}
*/

func (rt *_router) getMyGroups(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (rt *_router) deleteConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}

/*
func (rt *_router) getConversationMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}

*/

func (rt *_router) getMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}

/*
	func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusNotImplemented)
	}

func (rt *_router) listGroups(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (rt *_router) getGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (rt *_router) deleteGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (rt *_router) listGroupMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}


func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
}

*/
