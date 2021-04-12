package notes

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/theartofdevel/notes_system/api_service/internal/apperror"
	"github.com/theartofdevel/notes_system/api_service/internal/client/note_service"
	"github.com/theartofdevel/notes_system/api_service/pkg/jwt"
	"github.com/theartofdevel/notes_system/api_service/pkg/logging"
	"net/http"
)

const (
	notesURL = "/api/notes"
	noteURL  = "/api/notes/:uuid"
)

type Handler struct {
	Logger      logging.Logger
	NoteService note_service.NoteService
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, notesURL, jwt.Middleware(apperror.Middleware(h.GetNotes)))
	router.HandlerFunc(http.MethodPost, notesURL, jwt.Middleware(apperror.Middleware(h.CreateNote)))
	router.HandlerFunc(http.MethodGet, noteURL, jwt.Middleware(apperror.Middleware(h.GetNoteByUuid)))
	router.HandlerFunc(http.MethodPatch, noteURL, jwt.Middleware(apperror.Middleware(h.PartiallyUpdateNote)))
	router.HandlerFunc(http.MethodDelete, noteURL, jwt.Middleware(apperror.Middleware(h.DeleteNote)))
}

func (h *Handler) GetNotes(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	categoryUUID := r.URL.Query().Get("category_uuid")
	notes, err := h.NoteService.GetByCategoryUUID(r.Context(), categoryUUID)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(notes)

	return nil
}

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()
	var crNote note_service.CreateNoteDTO
	if err := json.NewDecoder(r.Body).Decode(&crNote); err != nil {
		return apperror.BadRequestError("can't decode")
	}

	noteUUID, err := h.NoteService.Create(r.Context(), crNote)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", fmt.Sprintf("%s/%s", notesURL, noteUUID))

	return nil
}

func (h *Handler) GetNoteByUuid(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	noteUuid := params.ByName("uuid")

	note, err := h.NoteService.GetByUUID(r.Context(), noteUuid)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(note)

	return nil
}

func (h *Handler) PartiallyUpdateNote(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	noteUUID := params.ByName("uuid")

	var dto note_service.UpdateNoteDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperror.BadRequestError("can't decode")
	}
	if err := h.NoteService.Update(r.Context(), noteUUID, dto); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	noteUUID := params.ByName("uuid")
	if err := h.NoteService.Delete(r.Context(), noteUUID); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
