package auth

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/theartofdevel/notes_system/api_service/internal/apperror"
	"github.com/theartofdevel/notes_system/api_service/internal/client/user_service"
	"github.com/theartofdevel/notes_system/api_service/pkg/jwt"
	"github.com/theartofdevel/notes_system/api_service/pkg/logging"
	"net/http"
)

const (
	authURL   = "/api/auth"
	signupURL = "/api/signup"
)

type Handler struct {
	Logger      logging.Logger
	UserService user_service.UserService
	JWTHelper   jwt.Helper
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, authURL, apperror.Middleware(h.Auth))
	router.HandlerFunc(http.MethodPut, authURL, apperror.Middleware(h.Auth))
	router.HandlerFunc(http.MethodPost, signupURL, apperror.Middleware(h.Signup))
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()
	var dto user_service.CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperror.BadRequestError("failed to decode data")
	}

	u, err := h.UserService.Create(r.Context(), dto)
	if err != nil {
		return err
	}
	token, err := h.JWTHelper.GenerateAccessToken(u)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(token)

	return nil
}

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	var token []byte
	var err error
	switch r.Method {
	case http.MethodPost:
		defer r.Body.Close()
		var dto user_service.SigninUserDTO
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			return apperror.BadRequestError("failed to decode data")
		}
		u, err := h.UserService.GetByEmailAndPassword(r.Context(), dto.Email, dto.Password)
		if err != nil {
			return err
		}
		token, err = h.JWTHelper.GenerateAccessToken(u)
		if err != nil {
			return err
		}
	case http.MethodPut:
		defer r.Body.Close()
		var rt jwt.RT
		if err := json.NewDecoder(r.Body).Decode(&rt); err != nil {
			return apperror.BadRequestError("failed to decode data")
		}
		token, err = h.JWTHelper.UpdateRefreshToken(rt)
		if err != nil {
			return err
		}
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(token)

	return err
}
