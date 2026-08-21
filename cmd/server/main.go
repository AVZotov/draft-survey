package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	draftsurvey "github.com/AVZotov/draft-survey"
	"github.com/AVZotov/draft-survey/internal/config"
	"github.com/AVZotov/draft-survey/internal/handler"
	"github.com/AVZotov/draft-survey/internal/logger"
	"github.com/AVZotov/draft-survey/internal/service"
	"github.com/AVZotov/draft-survey/internal/sse"
	"github.com/AVZotov/draft-survey/internal/storage"
	"github.com/AVZotov/draft-survey/internal/validation"
	"github.com/go-chi/chi/v5"
)

// findFreePort returns preferred if it can be bound, otherwise scans
// upward from preferred's own port number (falling back to 3399 only if
// preferred can't be parsed) until it finds one that binds, or gives up
// and returns preferred unchanged after 100 attempts.
func findFreePort(preferred string) string {
	ln, err := net.Listen("tcp", preferred)
	if err == nil {
		ln.Close()
		return preferred
	}

	base := 3399
	if port, perr := strconv.Atoi(strings.TrimPrefix(preferred, ":")); perr == nil {
		base = port
	}

	for port := base + 1; port < base+100; port++ {
		addr := fmt.Sprintf(":%d", port)
		ln, err := net.Listen("tcp", addr)
		if err == nil {
			ln.Close()
			return addr
		}
	}
	return preferred // fallback
}

func main() {
	cfg := config.Load()

	db, err := storage.NewDB(cfg.DBPath, draftsurvey.Dictionaries)
	if err != nil {
		log.Fatal(err)
	}

	slog := logger.NewSlogLogger(cfg.LogLevel)
	validator := validation.New()

	userStore := storage.NewSQLiteUserStore(db)
	dictionaryStore := storage.NewDictionariesStore(db)
	surveyStore := storage.NewSQLiteSurveyStore(db)

	services := &service.Services{
		User:       service.NewUserService(userStore, slog, validator),
		Survey:     service.NewSurveyService(surveyStore, userStore, slog),
		Draft:      service.NewDraftService(surveyStore, userStore, slog, validation.NewDraftValidator()),
		Tank:       service.NewTankService(surveyStore, slog),
		Dictionary: service.NewDictionaryService(dictionaryStore),
	}

	broker := sse.NewBroker()
	chiH := handler.New(services, slog, cfg.Version, broker)

	r := chi.NewRouter()
	if err := handler.SetupRoutesChi(r, chiH); err != nil {
		log.Fatal(err)
	}

	addr := findFreePort(cfg.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("forced shutdown", err)
	}

	if err = db.Close(); err != nil {
		slog.Error("error closing database", err)
	}
	slog.Info("server stopped")
}
