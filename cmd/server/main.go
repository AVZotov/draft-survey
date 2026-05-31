package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	draftsurvey "github.com/AVZotov/draft-survey"
	"github.com/AVZotov/draft-survey/internal/handler"
	chihandler "github.com/AVZotov/draft-survey/internal/handler/chi"
	"github.com/AVZotov/draft-survey/internal/logger"
	"github.com/AVZotov/draft-survey/internal/service"
	"github.com/AVZotov/draft-survey/internal/storage"
	"github.com/AVZotov/draft-survey/internal/validation"
	"github.com/go-chi/chi/v5"
	"github.com/gofiber/fiber/v2"
)

var version = "dev"

const DBPath = "./data/draft-survey.db"

func main() {
	db, err := storage.NewDB(DBPath)
	if err != nil {
		log.Fatal(err)
	}

	surveyStore := storage.NewSQLiteSurveyStore(db)
	userStore := storage.NewSQLiteUserStore(db)
	dictionariesStore := storage.NewDictionariesStore(draftsurvey.Dictionaries)

	// --- Old Fiber server (port 3399) ---
	app := fiber.New()
	validator := validation.New()
	h := handler.New(userStore, surveyStore, surveyStore, dictionariesStore, validator, version)
	handler.SetupRoutes(app, h)

	// --- New Chi server (port 3400) ---
	slog := logger.NewSlogLogger()
	services := &service.Services{
		User:       service.NewUserService(userStore, slog, validator),
		Survey:     &service.NoopSurveyService{},
		Draft:      &service.NoopDraftService{},
		Dictionary: &service.NoopDictionaryService{},
	}
	chiH := chihandler.New(services, &logger.Noop{}, version) // nil logger — NoopLogger later
	r := chi.NewRouter()
	if err := chihandler.SetupRoutesChi(r, chiH); err != nil {
		log.Fatal(err)
	}

	// Start both servers
	// go func() {
	// 	log.Fatal(app.Listen(":3399"))
	// }()

	go func() {
		log.Println("Chi server starting on :3400")
		log.Fatal(http.ListenAndServe(":3400", r))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	log.Println("Shutting down...")

	if err = app.Shutdown(); err != nil {
		log.Fatalf("Fiber shutdown error: %v", err)
	}
	if err = db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	log.Println("Server exiting")
}
