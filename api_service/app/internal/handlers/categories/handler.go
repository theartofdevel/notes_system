package categories

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/theartofdevel/notes_system/api_service/pkg/logging"
	"github.com/theartofdevel/notes_system/api_service/pkg/middleware/jwt"
	"net/http"
)

const (
	categoriesURL = "/api/categories"
	categoryURL = "/api/categories/:id"
)

type category struct {
	Id       string `json:"id"`
	ParentId string `json:"parent_id,omitempty"`
	Name     string `json:"name"`
	Color    string `json:"color"`
}

type createCategory struct {
	category
	UserId string `json:"user_id,omitempty"`
}

type Handler struct {
	Logger logging.Logger
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, categoriesURL, jwt.JWTMiddleware(h.GetCategories))
	router.HandlerFunc(http.MethodPost, categoriesURL, jwt.JWTMiddleware(h.CreateCategory))
	router.HandlerFunc(http.MethodPatch, categoryURL, jwt.JWTMiddleware(h.PartiallyUpdateCategory))
	router.HandlerFunc(http.MethodDelete, categoryURL, jwt.JWTMiddleware(h.DeleteCategory))
}

func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)
	// TODO call CategoryService
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(userId))
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)
	var crCategory createCategory
	if err := json.NewDecoder(r.Body).Decode(&crCategory); err != nil {
		h.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	crCategory.UserId = userId
	// TODO call CategoryService
	w.Header().Set("Content-Type", "application/json")
	// TODO set real category id
	w.Header().Set("Location", fmt.Sprintf("%s/%s", categoriesURL, "category_id"))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("to be done"))
}

func (h *Handler) PartiallyUpdateCategory(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	categoryId := params.ByName("id")
	var c category
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		h.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	c.Id = categoryId
	// TODO call CategoryService
	w.WriteHeader(http.StatusNoContent)
	// del it
	w.Write([]byte(userId))
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	categoryId := params.ByName("id")
	// TODO call CategoryService
	w.WriteHeader(http.StatusNoContent)
	// del it
	w.Write([]byte(userId))
	w.Write([]byte(categoryId))
}