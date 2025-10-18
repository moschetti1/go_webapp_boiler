package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/moschetti1/cronsearch/internal/config"
	"github.com/moschetti1/cronsearch/internal/handlers"
	custom_middleware "github.com/moschetti1/cronsearch/internal/middleware"
	"github.com/moschetti1/cronsearch/internal/repository"
	"github.com/moschetti1/cronsearch/internal/session"
	"github.com/moschetti1/cronsearch/web"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	config := config.LoadConfig()
	db, err := sql.Open("libsql", config.DBURL)

	if err != nil {
		log.Fatalf("Connection to db %s failed: %s", config.DBURL, err.Error())
	}
	repo := repository.New(db)

	sessionManager := session.NewSessionManager([]byte(config.CookieSecret))

	authHandler := handlers.NewAuthHandler(repo, sessionManager)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Handle("/static/*", http.FileServerFS(web.StaticFilesFS))

	router.Route("/auth/google", func(r chi.Router) {
		r.Get("/login", authHandler.GoogleLogin)
		r.Get("/callback", authHandler.GoogleCallback)
		r.Get("/logout", authHandler.Logout)
	})

	router.Route("/dashboard", func(r chi.Router) {
		r.Use(custom_middleware.UserCtx(repo, sessionManager))
	})

	http.ListenAndServe(config.Port, router)
}
