package auth

import (
	"encoding/json"
	"github.com/cristalhq/jwt/v3"
	"github.com/theartofdevel/notes_system/api_service/internal/config"
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

func Auth(w http.ResponseWriter, r *http.Request) {
	var u user
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		logging.GetLogger().Fatal(err)
	}

	defer r.Body.Close()
	// TODO client to UserService and get user by username and password
	// for now sub check
	if u.Username != "me" || u.Password != "pass" {
		w.WriteHeader(http.StatusNotFound)
		return
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
		logging.GetLogger().Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	jsonBytes, err := json.Marshal(map[string]string{
		"token": token.String(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
	w.WriteHeader(200)
}
