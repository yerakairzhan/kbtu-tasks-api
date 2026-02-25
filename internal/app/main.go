package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tasks_assignment/internal/handlers"
	"tasks_assignment/internal/logger"
	"tasks_assignment/internal/middleware"
	"tasks_assignment/internal/repository"
	"tasks_assignment/internal/repository/_postgres"
	"tasks_assignment/internal/usecase"
	"tasks_assignment/pkg/modules"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := modules.LoadConfig("config.yml")
	if err != nil {
		logger.Fatalf("load config: %v", err)
	}

	postgre := _postgres.NewPGXDialect(ctx, &cfg.PG)
	defer postgre.DB.Close()
	logger.Infof("database connection initialized")

	repos := repository.NewRepositories(postgre)
	taskUsecase := usecase.NewTaskUsecase(repos.Tasks)
	taskHandler := handlers.NewTaskHandler(taskUsecase)

	publicMux := http.NewServeMux()
	privateMux := http.NewServeMux()

	privateMux.HandleFunc("/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetTasks(w, r)
		case http.MethodPost:
			taskHandler.CreateTask(w, r)
		case http.MethodPatch:
			taskHandler.UpdateTask(w, r)
		case http.MethodDelete:
			taskHandler.DeleteTask(w, r)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	privateMux.HandleFunc("/v1/external-tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			taskHandler.FetchExternalTasks(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	publicMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	publicMux.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "docs/swagger.yaml")
	})

	publicMux.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(swaggerHTML))
	})

	privateHandler := middleware.APIKeyAuth(cfg.APIKey, privateMux)
	handler := middleware.Logging(middleware.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/healthz", "/swagger", "/swagger.yaml":
			publicMux.ServeHTTP(w, r)
		default:
			privateHandler.ServeHTTP(w, r)
		}
	})))

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}

	go func() {
		logger.Infof("server starting on :%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Infof("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatalf("server shutdown failed: %v", err)
	}

	logger.Infof("server stopped")
}

const swaggerHTML = `<!doctype html>
<html>
  <head>
    <meta charset="UTF-8" />
    <title>Tasks API Swagger</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.ui = SwaggerUIBundle({
        url: '/swagger.yaml',
        dom_id: '#swagger-ui'
      });
    </script>
  </body>
</html>`
