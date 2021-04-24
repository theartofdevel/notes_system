package note

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/theartofdevel/notes_system/note_service/internal/apperror"
	"github.com/theartofdevel/notes_system/note_service/pkg/logging"
	"net/http"
)

const (
	notesURL = "/api/notes"
	noteURL  = "/api/notes/:uuid"
)

type Handler struct {
	Logger      logging.Logger
	NoteService Service
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, noteURL, apperror.Middleware(h.GetNote))
	router.HandlerFunc(http.MethodGet, notesURL, apperror.Middleware(h.GetNotesByCategory))
	router.HandlerFunc(http.MethodPost, notesURL, apperror.Middleware(h.CreateNote))
	router.HandlerFunc(http.MethodPatch, noteURL, apperror.Middleware(h.PartiallyUpdateNote))
	router.HandlerFunc(http.MethodDelete, noteURL, apperror.Middleware(h.DeleteNote))
}

func (h *Handler) GetNote(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET NOTE")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	noteUUID := params.ByName("uuid")
	if noteUUID == "" {
		return apperror.BadRequestError("uuid query parameter is required and must be a comma separated integers")
	}

	note, err := h.NoteService.GetOne(r.Context(), noteUUID)
	if err != nil {
		return err
	}
	noteBytes, err := json.Marshal(note)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(noteBytes)

	return nil
}

func (h *Handler) GetNotesByCategory(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET NOTES BY CATEGORY")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get category_uuid from URL")
	categoryUUID := r.URL.Query().Get("category_uuid")
	if categoryUUID == "" {
		return apperror.BadRequestError("category_uuid query parameter is required and must be a comma separated integers")
	}

	notes, err := h.NoteService.GetByCategoryUUID(r.Context(), categoryUUID)
	if err != nil {
		return err
	}

	notesBytes, err := json.Marshal(notes)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(notesBytes)

	return nil
}

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("CREATE NOTE")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("decode create tag dto")
	var dto CreateNoteDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperror.BadRequestError("invalid data")
	}

	noteUUID, err := h.NoteService.Create(r.Context(), dto)
	if err != nil {
		return err
	}
	w.Header().Set("Location", fmt.Sprintf("%s/%s", notesURL, noteUUID))
	w.WriteHeader(http.StatusCreated)

	return nil
}

func (h *Handler) PartiallyUpdateNote(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("PARTIALLY UPDATE NOTE")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	noteUUID := params.ByName("uuid")
	if noteUUID == "" {
		return apperror.BadRequestError("uuid query parameter is required and must be a comma separated integers")
	}

	h.Logger.Debug("decode update tag dto")
	var dto UpdateNoteDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperror.BadRequestError("invalid data")
	}

	dto.UUID = noteUUID

	err := h.NoteService.Update(r.Context(), dto)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("DELETE NOTE")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	noteUUID := params.ByName("uuid")
	if noteUUID == "" {
		return apperror.BadRequestError("uuid query parameter is required and must be a comma separated integers")
	}

	err := h.NoteService.Delete(r.Context(), noteUUID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
