package user

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/theartofdevel/notes_system/user_service/internal/apperror"
	"github.com/theartofdevel/notes_system/user_service/pkg/logging"
	"net/http"
)

const (
	usersURL = "/api/users"
	userURL  = "/api/users/:uuid"
)

type Handler struct {
	Logger      logging.Logger
	UserService Service
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, usersURL, apperror.Middleware(h.GetUserByEmailAndPassword))
	router.HandlerFunc(http.MethodPost, usersURL, apperror.Middleware(h.CreateUser))
	router.HandlerFunc(http.MethodGet, userURL, apperror.Middleware(h.GetUser))
	router.HandlerFunc(http.MethodPatch, userURL, apperror.Middleware(h.PartiallyUpdateUser))
	router.HandlerFunc(http.MethodDelete, userURL, apperror.Middleware(h.DeleteUser))
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET USER")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userUUID := params.ByName("uuid")

	user, err := h.UserService.GetOne(r.Context(), userUUID)
	if err != nil {
		return err
	}

	h.Logger.Debug("marshal user")
	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshall user. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}

func (h *Handler) GetUserByEmailAndPassword(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET USER BY EMAIL AND PASSWORD")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get email and password from URL")
	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")
	if email == "" || password == "" {
		return apperror.BadRequestError("invalid query parameters email or password")
	}

	user, err := h.UserService.GetByEmailAndPassword(r.Context(), email, password)
	if err != nil {
		return err
	}

	h.Logger.Debug("marshal user")
	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)

	return nil
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("CREATE USER")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("decode create user dto")
	var crUser CreateUserDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&crUser); err != nil {
		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
	}

	userUUID, err := h.UserService.Create(r.Context(), crUser)
	if err != nil {
		return err
	}
	w.Header().Set("Location", fmt.Sprintf("%s/%s", usersURL, userUUID))
	w.WriteHeader(http.StatusCreated)

	return nil
}

func (h *Handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("PARTIALLY UPDATE USER")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userUUID := params.ByName("uuid")

	h.Logger.Debug("decode update user dto")
	var updUser UpdateUserDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&updUser); err != nil {
		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
	}
	updUser.UUID = userUUID

	err := h.UserService.Update(r.Context(), updUser)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("DELETE USER")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userUUID := params.ByName("uuid")

	err := h.UserService.Delete(r.Context(), userUUID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
