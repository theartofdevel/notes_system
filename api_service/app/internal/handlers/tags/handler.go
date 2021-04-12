package tags

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/theartofdevel/notes_system/api_service/internal/apperror"
	"github.com/theartofdevel/notes_system/api_service/internal/client/tag_service"
	"github.com/theartofdevel/notes_system/api_service/pkg/jwt"
	"github.com/theartofdevel/notes_system/api_service/pkg/logging"
	"net/http"
	"strconv"
	"strings"
)

const (
	tagsURL = "/api/tags"
	tagURL  = "/api/tags/:id"
)

type Handler struct {
	Logger     logging.Logger
	TagService tag_service.TagService
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, tagURL, jwt.Middleware(apperror.Middleware(h.GetTag)))
	router.HandlerFunc(http.MethodGet, tagsURL, jwt.Middleware(apperror.Middleware(h.GetManyTags)))
	router.HandlerFunc(http.MethodPost, tagsURL, jwt.Middleware(apperror.Middleware(h.CreateTag)))
	router.HandlerFunc(http.MethodPatch, tagURL, jwt.Middleware(apperror.Middleware(h.PartiallyUpdateTag)))
	router.HandlerFunc(http.MethodDelete, tagURL, jwt.Middleware(apperror.Middleware(h.DeleteTag)))
}

func (h *Handler) GetTag(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	tagIDStr := params.ByName("id")
	id, err := strconv.Atoi(tagIDStr)
	if err != nil {
		return apperror.BadRequestError("invalid id")
	}

	tag, err := h.TagService.GetOne(r.Context(), id)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(tag)

	return nil
}

func (h *Handler) GetManyTags(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	idsParam := r.URL.Query().Get("id")
	if idsParam == "" {
		return apperror.BadRequestError("invalid id")
	}

	var tagsIds []int
	idsStr := strings.Split(idsParam, ",")
	for _, idStr := range idsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return apperror.BadRequestError("invalid id")
		}
		tagsIds = append(tagsIds, id)
	}

	tags, err := h.TagService.GetMany(r.Context(), tagsIds)
	if err != nil {
		return err
	}
	if len(tags) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return apperror.ErrNotFound
	}

	w.WriteHeader(http.StatusOK)
	w.Write(tags)

	return nil
}

func (h *Handler) CreateTag(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	if r.Context().Value("user_uuid") == nil {
		h.Logger.Error("there is no user_uuid in context")
		return apperror.UnauthorizedError("")
	}
	userUUID := r.Context().Value("user_uuid").(string)

	var dto tag_service.CreateTagDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperror.BadRequestError("can't decode")
	}
	dto.UserUUID = userUUID

	tagID, err := h.TagService.Create(r.Context(), dto)
	if err != nil {
		return err
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%s", tagsURL, tagID))
	w.WriteHeader(http.StatusCreated)

	return nil
}

func (h *Handler) PartiallyUpdateTag(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	if r.Context().Value("user_uuid") == nil {
		h.Logger.Error("there is no user_uuid in context")
		return apperror.UnauthorizedError("")
	}
	userUUID := r.Context().Value("user_uuid").(string)

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	tagId := params.ByName("id")

	var dto tag_service.UpdateTagDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperror.BadRequestError("can't decode")
	}
	dto.UserUUID = userUUID
	if err := h.TagService.Update(r.Context(), tagId, dto); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *Handler) DeleteTag(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	if r.Context().Value("user_uuid") == nil {
		h.Logger.Error("there is no user_uuid in context")
		return apperror.UnauthorizedError("")
	}

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	tagId := params.ByName("id")
	if err := h.TagService.Delete(r.Context(), tagId); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
