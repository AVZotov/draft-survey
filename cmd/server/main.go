package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	
	draftsurvey "github.com/AVZotov/draft-survey"
	"github.com/AVZotov/draft-survey/internal/handler"
	"github.com/AVZotov/draft-survey/internal/logger"
	"github.com/AVZotov/draft-survey/internal/service"
	"github.com/AVZotov/draft-survey/internal/sse"
	"github.com/AVZotov/draft-survey/internal/storage"
	"github.com/AVZotov/draft-survey/internal/validation"
	"github.com/go-chi/chi/v5"
)

var version = "dev"

const DBPath = "./data/draft-survey.db"

func main() {
	db, err := storage.NewDB(DBPath, draftsurvey.Dictionaries)
	if err != nil {
		log.Fatal(err)
	}
	
	slog := logger.NewSlogLogger()
	validator := validation.New()
	
	userStore := storage.NewSQLiteUserStore(db)
	dictionaryStore := storage.NewDictionariesStore(db)
	
	services := &service.Services{
		User:       service.NewUserService(userStore, slog, validator),
		Survey:     &service.NoopSurveyService{},
		Draft:      &service.NoopDraftService{},
		Dictionary: service.NewDictionaryService(dictionaryStore),
	}
	
	broker := sse.NewBroker()
	chiH := handler.New(services, slog, version, broker)
	
	r := chi.NewRouter()
	if err := handler.SetupRoutesChi(r, chiH); err != nil {
		log.Fatal(err)
	}
	
	go func() {
		log.Println("Server starting on :3400")
		log.Fatal(http.ListenAndServe(":3400", r))
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
