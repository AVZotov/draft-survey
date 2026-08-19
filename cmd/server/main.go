package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	go func() {
		log.Printf("Server starting on %s", cfg.Port)
		log.Fatal(http.ListenAndServe(cfg.Port, r))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down...")
	if err = db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
	log.Println("Server exiting")
}
