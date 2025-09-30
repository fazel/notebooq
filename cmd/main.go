package main

import (
	"log"
	"net/http"
	"time"

	"github.com/fazel/notebooq/internal/auth"
	"github.com/fazel/notebooq/internal/config"
	"github.com/fazel/notebooq/internal/db"
	"github.com/fazel/notebooq/internal/handlers"
	"github.com/fazel/notebooq/internal/models"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dbConn, err := db.Open(cfg)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	// run migrations
	dbConn.AutoMigrate(&models.User{}, &models.Note{})

	r := gin.Default()

	// public routes
	r.POST("/api/signup", handlers.NewUserHandler(dbConn, cfg).Signup)
	r.POST("/api/login", handlers.NewUserHandler(dbConn, cfg).Login)

	// protected
	authMiddleware := auth.JWTMiddleware(cfg.JWTSecret)
	authGroup := r.Group("/api")
	authGroup.Use(authMiddleware)
	{
		authGroup.GET("/notes", handlers.NewNoteHandler(dbConn).ListNotes)
		authGroup.POST("/notes", handlers.NewNoteHandler(dbConn).CreateNote)
		authGroup.GET("/notes/:id", handlers.NewNoteHandler(dbConn).GetNote)
		authGroup.PUT("/notes/:id", handlers.NewNoteHandler(dbConn).UpdateNote)
		authGroup.DELETE("/notes/:id", handlers.NewNoteHandler(dbConn).DeleteNote)
	}

	srv := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("listening on %s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
