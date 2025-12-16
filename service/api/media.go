package api

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// uploadImage gestisce l'upload di un file immagine (POST /images)
func (rt *_router) uploadImage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 1. Controllo Autenticazione (Solo utenti loggati possono caricare)
	_, err := rt.extractUser(r)
	if err != nil {
		rt.sendJSONError(w, http.StatusUnauthorized, "Unauthorized: please login first")
		return
	}

	// 2. Parsing del form multipart (Limite 10MB)
	const maxUploadSize = 10 << 20 // 10 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		rt.sendJSONError(w, http.StatusBadRequest, "File too big (max 10MB) or invalid format")
		return
	}

	// 3. Recupera il file dal form (chiave "file" come da YAML)
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		rt.sendJSONError(w, http.StatusBadRequest, "Missing 'file' field in multipart/form-data")
		return
	}
	defer func() {
		_ = file.Close()
	}()
	// 4. Lettura primi bytes per controllare il tipo (MIME sniffing)
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		rt.sendJSONError(w, http.StatusInternalServerError, "Error reading file")
		return
	}

	fileType := http.DetectContentType(buff)
	if fileType != "image/jpeg" && fileType != "image/png" {
		rt.sendJSONError(w, http.StatusBadRequest, "Only .jpg and .png images are allowed")
		return
	}

	// Reset del puntatore del file all'inizio dopo la lettura
	_, _ = file.Seek(0, 0)

	// 5. Generazione nome univoco (UUID)
	id := uuid.New().String()
	fileExt := filepath.Ext(fileHeader.Filename)
	if fileExt == "" {
		// Fallback se non c'è estensione
		if fileType == "image/png" {
			fileExt = ".png"
		} else {
			fileExt = ".jpg"
		}
	}
	newFileName := id + fileExt

	// Percorso di salvataggio (cartella "images" nella root)
	savePath := filepath.Join("./images", newFileName)

	// 6. Salvataggio su disco
	dst, err := os.Create(savePath)
	if err != nil {

		if os.IsNotExist(err) {
			_ = os.MkdirAll("./images", 0755) // Prova a creare la cartella
			dst, err = os.Create(savePath)    // Riprova
		}

		if err != nil {
			rt.baseLogger.WithError(err).Error("Error creating file on disk")
			rt.sendJSONError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}
	defer func() {
		_ = dst.Close()
	}()
	if _, err := io.Copy(dst, file); err != nil {
		rt.baseLogger.WithError(err).Error("Error saving file content")
		rt.sendJSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// 7. Costruzione Risposta
	publicURL := "/images/" + newFileName

	resp := ImageMetadata{
		ID:          id,
		URL:         publicURL,
		ContentType: fileType,
	}

	rt.sendJSONResponse(w, http.StatusCreated, resp)
}
