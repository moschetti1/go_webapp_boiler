package session

import "github.com/gorilla/sessions"

type SessionManager struct {
	Store *sessions.CookieStore
}

func NewSessionManager(secret []byte) *SessionManager {
	store := sessions.NewCookieStore(secret)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	return &SessionManager{Store: store}
}
