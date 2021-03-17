package notes

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/theartofdevel/notes_system/api_service/pkg/logging"
	"github.com/theartofdevel/notes_system/api_service/pkg/middleware/jwt"
	"net/http"
)

const (
	notesURL = "/api/notes"
	noteURL = "/api/notes/:uuid"
)

type note struct {
	Uuid        string `json:"uuid"`
	Header      string `json:"header"`
	Body        string `json:"body"`
	CreatedDate int    `json:"created_date,omitempty"`
	CategoryId  string `json:"category_id"`
}

type createNote struct {
	note
	UserId string `json:"user_id,omitempty"`
}

type Handler struct {
	Logger logging.Logger
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, notesURL, jwt.JWTMiddleware(h.GetNotes))
	router.HandlerFunc(http.MethodPost, notesURL, jwt.JWTMiddleware(h.CreateNote))
	router.HandlerFunc(http.MethodGet, noteURL, jwt.JWTMiddleware(h.GetNoteByUuid))
	router.HandlerFunc(http.MethodPatch, noteURL, jwt.JWTMiddleware(h.PartiallyUpdateNote))
	router.HandlerFunc(http.MethodDelete, noteURL, jwt.JWTMiddleware(h.DeleteNote))
}

func (h *Handler) GetNotes(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)
	// TODO call NoteService
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(userId))
}

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)
	var crNote createNote
	if err := json.NewDecoder(r.Body).Decode(&crNote); err != nil {
		h.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	crNote.UserId = userId
	// TODO call NoteService
	w.WriteHeader(http.StatusCreated)
	// TODO set real note uuid
	w.Header().Set("Location", fmt.Sprintf("%s/%s", notesURL, "note_uuid"))
}

func (h *Handler) GetNoteByUuid(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)
	// TODO call NoteService
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(userId))
}

func (h *Handler) PartiallyUpdateNote(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	noteUuid := params.ByName("uuid")
	var n note
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		h.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	n.Uuid = noteUuid
	// TODO call NoteService
	w.WriteHeader(http.StatusNoContent)
	// del it
	w.Write([]byte(userId))
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	noteUuid := params.ByName("uuid")
	// TODO call NoteService
	w.WriteHeader(http.StatusNoContent)
	// del it
	w.Write([]byte(userId))
	w.Write([]byte(noteUuid))
}