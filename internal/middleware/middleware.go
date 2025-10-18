package custom_middleware

import (
	"context"
	"net/http"

	"github.com/moschetti1/cronsearch/internal/repository"
	"github.com/moschetti1/cronsearch/internal/session"
)

var userContextKey = "user"

func UserCtx(q *repository.Queries, sm *session.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess, _ := sm.Store.Get(r, "session")
			id, ok := sess.Values["user_id"]
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			user, err := q.GetUserByGoogleId(r.Context(), id.(string))
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserFromCtx(r *http.Request) (repository.User, bool) {
	u, ok := r.Context().Value(userContextKey).(repository.User)
	return u, ok
}
