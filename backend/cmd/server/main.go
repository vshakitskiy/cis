package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/vshakitskiy/cis/internal/config"
	"github.com/vshakitskiy/cis/internal/db"
	"github.com/vshakitskiy/cis/internal/handler"
	"github.com/vshakitskiy/cis/internal/middleware"
	"github.com/vshakitskiy/cis/internal/repository"
	"github.com/vshakitskiy/cis/internal/response"
	"github.com/vshakitskiy/cis/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	ctx := context.Background()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connecting to database: %v", err)
	}
	defer pool.Close()

	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, response.JSON{"status": "ok"})
	})

	userRepo := repository.NewUserRepo(pool)
	projectRepo := repository.NewProjectRepo(pool)
	taskRepo := repository.NewTaskRepo(pool)
	timeEntryRepo := repository.NewTimeEntryRepo(pool)
	commentRepo := repository.NewCommentRepo(pool)
	attachmentRepo := repository.NewAttachmentRepo(pool)
	reportRepo := repository.NewReportRepo(pool)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	projectService := service.NewProjectService(projectRepo)
	taskService := service.NewTaskService(taskRepo, projectRepo)
	timeEntryService := service.NewTimeEntryService(timeEntryRepo, taskRepo)
	commentService := service.NewCommentService(commentRepo, taskRepo)
	attachmentService := service.NewAttachmentService(attachmentRepo, taskRepo, "./uploads")

	authHandler := handler.NewAuthHandler(authService)
	timeEntryHandler := handler.NewTimeEntryHandler(timeEntryService)
	commentHandler := handler.NewCommentHandler(commentService)
	attachmentHandler := handler.NewAttachmentHandler(attachmentService)
	taskHandler := handler.NewTaskHandler(taskService, timeEntryHandler, commentHandler, attachmentHandler)
	reportHandler := handler.NewReportHandler(reportRepo)
	projectHandler := handler.NewProjectHandler(projectService, taskHandler, reportHandler)

	authMiddleware := middleware.Auth(authService, userRepo)

	r.Route("/api", func(r chi.Router) {
		r.Mount("/auth", authHandler.Routes())

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Mount("/projects", projectHandler.Routes())
		})
	})

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		log.Println("shutting down server...")
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("server shutdown: %v", err)
		}
	}()

	log.Printf("server starting on :%s", cfg.ServerPort)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
