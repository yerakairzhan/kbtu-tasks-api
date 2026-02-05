package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tasks_assignment/internal/handlers"
	"tasks_assignment/internal/middleware"
	"time"
)

func main() {
	store := handlers.NewTaskStore()

	http.HandleFunc("/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			store.GetTasks(w, r)
		case http.MethodPost:
			store.CreateTask(w, r)
		case http.MethodPatch:
			store.UpdateTask(w, r)
		case http.MethodDelete:
			store.DeleteTask(w, r)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/v1/external-tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			store.FetchExternalTasks(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	handler := middleware.Logging(middleware.RequestID(middleware.APIKeyAuth(http.DefaultServeMux)))

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped")
}
