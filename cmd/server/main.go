package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	draft_survey "github.com/AVZotov/draft-survey"
	"github.com/AVZotov/draft-survey/internal/handler"
	"github.com/AVZotov/draft-survey/internal/storage"
	"github.com/AVZotov/draft-survey/internal/validation"
	"github.com/gofiber/fiber/v2"
)

var version = "dev"

const DBPath = "./data/draft-survey.db"

func main() {
	app := fiber.New()

	db, err := storage.NewDB(DBPath)
	if err != nil {
		log.Fatal(err)
	}

	surveyStore := storage.NewSQLiteSurveyStore(db)
	userStore := storage.NewSQLiteUserStore(db)
	dictionariesStore := storage.NewDictionariesStore(draft_survey.Dictionaries)
	validator := validation.New()

	h := handler.New(userStore, surveyStore, surveyStore, dictionariesStore, validator, version)

	handler.SetupRoutes(app, h)

	go func() {
		log.Fatal(app.Listen(":3399"))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	log.Println("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	log.Println("Server exiting")
}
