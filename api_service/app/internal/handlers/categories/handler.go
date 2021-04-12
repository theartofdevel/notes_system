package categories

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/theartofdevel/notes_system/api_service/internal/apperror"
	"github.com/theartofdevel/notes_system/api_service/internal/client/category_service"
	"github.com/theartofdevel/notes_system/api_service/pkg/jwt"
	"github.com/theartofdevel/notes_system/api_service/pkg/logging"
	"net/http"
)

const (
	categoriesURL = "/api/categories"
	categoryURL   = "/api/categories/:uuid"
)

type Handler struct {
	CategoryService category_service.CategoryService
	Logger          logging.Logger
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, categoriesURL, jwt.Middleware(apperror.Middleware(h.GetCategories)))
	router.HandlerFunc(http.MethodPost, categoriesURL, jwt.Middleware(apperror.Middleware(h.CreateCategory)))
	router.HandlerFunc(http.MethodPatch, categoryURL, jwt.Middleware(apperror.Middleware(h.PartiallyUpdateCategory)))
	router.HandlerFunc(http.MethodDelete, categoryURL, jwt.Middleware(apperror.Middleware(h.DeleteCategory)))
}

func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	if r.Context().Value("user_uuid") == nil {
		h.Logger.Error("there is no user_uuid in context")
		return apperror.UnauthorizedError("")
	}
	userUuid := r.Context().Value("user_uuid").(string)
	categories, err := h.CategoryService.GetUserCategories(r.Context(), userUuid)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	w.Write(categories)

	return nil
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	if r.Context().Value("user_uuid") == nil {
		h.Logger.Error("there is no user_uuid in context")
		return apperror.UnauthorizedError("")
	}
	userUuid := r.Context().Value("user_uuid").(string)

	var crCategory category_service.CreateCategoryDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&crCategory); err != nil {
		return apperror.BadRequestError("can't decode")
	}
	crCategory.UserUuid = userUuid

	categoryUuid, err := h.CategoryService.CreateCategory(r.Context(), crCategory)
	if err != nil {
		return err
	}
	w.Header().Set("Location", fmt.Sprintf("%s/%s", categoriesURL, categoryUuid))
	w.WriteHeader(http.StatusCreated)

	return nil
}

func (h *Handler) PartiallyUpdateCategory(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	if r.Context().Value("user_uuid") == nil {
		h.Logger.Error("there is no user_uuid in context")
		return apperror.UnauthorizedError("")
	}
	userUuid := r.Context().Value("user_uuid").(string)

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	categoryUuid := params.ByName("uuid")
	var categoryDTO category_service.UpdateCategoryDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&categoryDTO); err != nil {
		return apperror.BadRequestError("can't decode")
	}
	categoryDTO.UserUuid = userUuid
	err := h.CategoryService.UpdateCategory(r.Context(), categoryUuid, categoryDTO)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	if r.Context().Value("user_uuid") == nil {
		h.Logger.Error("there is no user_uuid in context")
		return apperror.UnauthorizedError("")
	}

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	categoryDTO := category_service.DeleteCategoryDTO{
		Uuid:     params.ByName("uuid"),
		UserUuid: r.Context().Value("user_uuid").(string),
	}
	err := h.CategoryService.DeleteCategory(r.Context(), categoryDTO)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
