package tags

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/theartofdevel/notes_system/api_service/pkg/logging"
	"github.com/theartofdevel/notes_system/api_service/pkg/middleware/jwt"
	"net/http"
)

const (
	tagsURL = "/api/tags"
	tagURL = "/api/tags/:id"
)

type tag struct {
	Id       string `json:"id"`
	ParentId string `json:"parent_id,omitempty"`
	Name     string `json:"name"`
	Color    string `json:"color"`
}

type createTag struct {
	tag
	UserId string `json:"user_id,omitempty"`
}

type Handler struct {
	Logger logging.Logger
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, tagsURL, jwt.Middleware(h.CreateTag))
	router.HandlerFunc(http.MethodPatch, tagURL, jwt.Middleware(h.PartiallyUpdateTag))
	router.HandlerFunc(http.MethodDelete, tagURL, jwt.Middleware(h.DeleteTag))
}

func (h *Handler) CreateTag(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)
	var crTag createTag
	if err := json.NewDecoder(r.Body).Decode(&crTag); err != nil {
		h.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	crTag.UserId = userId
	// TODO call TagService
	w.Header().Set("Content-Type", "application/json")
	// TODO set real tag id
	w.Header().Set("Location", fmt.Sprintf("%s/%s", tagsURL, "tag_id"))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("to be done"))
}

func (h *Handler) PartiallyUpdateTag(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	tagId := params.ByName("id")
	var c tag
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		h.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	c.Id = tagId
	// TODO call TagService
	w.WriteHeader(http.StatusNoContent)
	// del it
	w.Write([]byte(userId))
}

func (h *Handler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	tagId := params.ByName("id")
	// TODO call TagService
	w.WriteHeader(http.StatusNoContent)
	// del it
	w.Write([]byte(userId))
	w.Write([]byte(tagId))
}