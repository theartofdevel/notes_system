package auth

import (
	"encoding/json"
	"github.com/cristalhq/jwt/v3"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/theartofdevel/notes_system/api_service/internal/config"
	"github.com/theartofdevel/notes_system/api_service/pkg/cache"
	"github.com/theartofdevel/notes_system/api_service/pkg/logging"
	jwt2 "github.com/theartofdevel/notes_system/api_service/pkg/middleware/jwt"
	"net/http"
	"time"
)

const (
	authURL = "/api/auth"
	signupURL = "/api/signup"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type newUser struct {
	user
	Email string `json:"email"`
}

type refresh struct {
	RefreshToken string `json:"refresh_token"`
}

type Handler struct {
	Logger logging.Logger
	RTCache cache.Repository
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, authURL, h.Auth)
	router.HandlerFunc(http.MethodPut, authURL, h.Auth)
	router.HandlerFunc(http.MethodPost, signupURL, h.Signup)
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var nu newUser
	if err := json.NewDecoder(r.Body).Decode(&nu); err != nil {
		h.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// TODO validate username and password
	// TODO create user using UserService
	jsonBytes, errCode := h.generateAccessToken()
	if errCode != 0 {
		w.WriteHeader(errCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(jsonBytes)
}

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var u user
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			h.Logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		defer r.Body.Close()
		// TODO client to UserService and get user by username and password
		// for now stub check
		if u.Username != "me" || u.Password != "pass" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	case http.MethodPut:
		var refreshTokenS refresh
		if err := json.NewDecoder(r.Body).Decode(&refreshTokenS); err != nil {
			h.Logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		userIdBytes, err := h.RTCache.Get([]byte(refreshTokenS.RefreshToken))
		h.Logger.Info("refresh token user_id: %s", userIdBytes)
		if err != nil {
			h.Logger.Error(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.RTCache.Del([]byte(refreshTokenS.RefreshToken))
		// TODO client to UserService and get user by username
	}

	jsonBytes, errCode := h.generateAccessToken()
	if errCode != 0 {
		w.WriteHeader(errCode)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(jsonBytes)
}

func (h *Handler) generateAccessToken() ([]byte, int) {
	key := []byte(config.GetConfig().JWT.Secret)
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return nil, 418
	}
	builder := jwt.NewBuilder(signer)

	// TODO insert real user data in claims
	claims := jwt2.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "uuid_here",
			Audience:  []string{"users"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
		Email: "email@will.be.here",
	}
	token, err := builder.Build(claims)
	if err != nil {
		h.Logger.Error(err)
		return nil, http.StatusUnauthorized
	}

	h.Logger.Info("create refresh token")
	refreshTokenUuid := uuid.New()
	err = h.RTCache.Set([]byte(refreshTokenUuid.String()), []byte("user_uuid"), 0)
	if err != nil {
		h.Logger.Error(err)
		return nil, http.StatusInternalServerError
	}

	jsonBytes, err := json.Marshal(map[string]string{
		"token": token.String(),
		"refresh_token": refreshTokenUuid.String(),
	})
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	return jsonBytes, 0
}
