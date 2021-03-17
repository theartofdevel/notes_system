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
	URL = "/api/auth"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type refresh struct {
	RefreshToken string `json:"refresh_token"`
}

type Handler struct {
	Logger logging.Logger
	RTCache cache.Repository
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, URL, h.Auth)
	router.HandlerFunc(http.MethodPut, URL, h.Auth)
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

	key := []byte(config.GetConfig().JWT.Secret)
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		w.WriteHeader(418)
		return
	}
	builder := jwt.NewBuilder(signer)

	// TODO insert real user data in claims
	claims := jwt2.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "uuid_here",
			Audience:  []string{"users"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 5)),
		},
		Email: "email@will.be.here",
	}
	token, err := builder.Build(claims)
	if err != nil {
		h.Logger.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	h.Logger.Info("create refresh token")
	refreshTokenUuid := uuid.New()
	err = h.RTCache.Set([]byte(refreshTokenUuid.String()), []byte("user_uuid"), 0)
	if err != nil {
		h.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(map[string]string{
		"token": token.String(),
		"refresh_token": refreshTokenUuid.String(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
	w.WriteHeader(200)
}
