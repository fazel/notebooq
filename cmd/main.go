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

	// migrations
	if err := dbConn.AutoMigrate(&models.User{}, &models.Note{}); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	userHandler := handlers.NewUserHandler(dbConn, cfg)
	noteHandler := handlers.NewNoteHandler(dbConn)

	r := gin.Default()

	// فایل‌های استاتیک frontend
	r.Static("/js", "./frontend/js")
	r.StaticFile("/login.html", "./frontend/login.html")
	r.StaticFile("/signup.html", "./frontend/signup.html")
	r.StaticFile("/notes.html", "./frontend/notes.html")
	r.StaticFile("/verify.html", "./frontend/verify.html")
	r.StaticFile("/profile.html", "./frontend/profile.html")

	// default / → redirect login
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/login.html")
	})

	// public API
	r.POST("/api/signup", userHandler.Signup)
	r.POST("/api/login", userHandler.Login)
	r.POST("/api/verify-email", userHandler.VerifyEmail)

	// protected API
	authMiddleware := auth.JWTMiddleware(cfg.JWTSecret)
	authGroup := r.Group("/api")
	authGroup.Use(authMiddleware)
	{
		authGroup.GET("/notes", noteHandler.ListNotes)
		authGroup.POST("/notes", noteHandler.CreateNote)
		authGroup.GET("/notes/:id", noteHandler.GetNote)
		authGroup.PUT("/notes/:id", noteHandler.UpdateNote)
		authGroup.DELETE("/notes/:id", noteHandler.DeleteNote)
		authGroup.GET("/profile", userHandler.Profile)
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
