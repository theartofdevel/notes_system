package categories

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/theartofdevel/notes_system/api_service/internal/client/category"
	"github.com/theartofdevel/notes_system/api_service/pkg/logging"
	"github.com/theartofdevel/notes_system/api_service/pkg/middleware/jwt"
	"net/http"
)

const (
	categoriesURL = "/api/categories"
	categoryURL   = "/api/categories/:uuid"
)

type Handler struct {
	CategoryService *category.Client
	Logger          logging.Logger
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, categoriesURL, jwt.JWTMiddleware(h.GetCategories))
	router.HandlerFunc(http.MethodPost, categoriesURL, jwt.JWTMiddleware(h.CreateCategory))
	router.HandlerFunc(http.MethodPatch, categoryURL, jwt.JWTMiddleware(h.PartiallyUpdateCategory))
	router.HandlerFunc(http.MethodDelete, categoryURL, jwt.JWTMiddleware(h.DeleteCategory))
}

func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("user_uuid") == nil {
		h.Logger.Error("there is no user_uuid in context")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userUuid := r.Context().Value("user_uuid").(string)
	categories, err := h.CategoryService.GetCategories(userUuid, context.Background(), nil)
	if err != nil {
		h.Logger.Error(err)
		w.WriteHeader(418)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(categories)
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("user_uuid") == nil {
		h.Logger.Error("there is no user_uuid in context")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userUuid := r.Context().Value("user_uuid").(string)

	var crCategory category.CreateCategoryDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&crCategory); err != nil {
		h.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	crCategory.UserUuid = userUuid
	w.Header().Set("Content-Type", "application/json")

	categoryUuid, err := h.CategoryService.CreateCategory(crCategory, context.Background(), nil)
	if err != nil {
		h.Logger.Error(err)
		w.WriteHeader(418)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s/%s", categoriesURL, categoryUuid))
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) PartiallyUpdateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("user_uuid") == nil {
		h.Logger.Error("there is no user_uuid in context")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userUuid := r.Context().Value("user_uuid").(string)

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	categoryUuid := params.ByName("uuid")
	var categoryDTO category.UpdateCategoryDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&categoryDTO); err != nil {
		h.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	categoryDTO.Uuid = categoryUuid
	categoryDTO.UserUuid = userUuid
	err := h.CategoryService.UpdateCategory(categoryDTO, context.Background(), nil)
	if err != nil {
		h.Logger.Error(err)
		w.WriteHeader(418)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("user_uuid") == nil {
		h.Logger.Error("there is no user_uuid in context")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	categoryDTO := category.DeleteCategoryDTO{
		Uuid:     params.ByName("uuid"),
		UserUuid: r.Context().Value("user_uuid").(string),
	}
	err := h.CategoryService.DeleteCategory(categoryDTO, context.Background(), nil)
	if err != nil {
		h.Logger.Error(err)
		w.WriteHeader(418)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
