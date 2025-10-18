package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/moschetti1/cronsearch/internal/config"
	"github.com/moschetti1/cronsearch/internal/repository"
	"github.com/moschetti1/cronsearch/internal/session"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
	Queries  *repository.Queries
	Sessions *session.SessionManager
	Google   *oauth2.Config
}

func NewAuthHandler(queries *repository.Queries, sessions *session.SessionManager) *AuthHandler {
	cfg := config.MustGet()
	googleClient := &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/auth/google/callback", cfg.PublicHost),
		ClientID:     cfg.GoogleClientId,
		ClientSecret: cfg.GoogleSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return &AuthHandler{
		Queries:  queries,
		Sessions: sessions,
		Google:   googleClient,
	}
}

func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := h.Google.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := h.Google.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Token exchange failed", http.StatusInternalServerError)
		return
	}
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var uinfo struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	json.NewDecoder(resp.Body).Decode((&uinfo))
	ctx := r.Context()
	user, err := h.Queries.GetUserByGoogleId(ctx, uinfo.ID)
	if err == sql.ErrNoRows {
		user, err = h.Queries.CreateUser(ctx, repository.CreateUserParams{
			ID:        uuid.New().String(),
			GoogleID:  uinfo.ID,
			Email:     uinfo.Email,
			Name:      sql.NullString{String: uinfo.Name, Valid: uinfo.Name != ""},
			AvatarUrl: sql.NullString{String: uinfo.Picture, Valid: uinfo.Picture != ""},
		})
		if err != nil {
			http.Error(w, "User creation failed", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	sess, _ := h.Sessions.Store.Get(r, "session")
	sess.Values["user_id"] = user.ID
	sess.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sess, _ := h.Sessions.Store.Get(r, "session")
	delete(sess.Values, "user_id")
	sess.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
