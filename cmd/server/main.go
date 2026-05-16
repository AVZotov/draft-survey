package main

import (
	"log"

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

	userStore, err := storage.NewUserStore("./data/users")
	if err != nil {
		log.Fatal(err)
	}

	surveyStore, err := storage.NewSurveyStore("./data/surveys", "./data/temp")
	if err != nil {
		log.Fatal(err)
	}

	surveyQueryStore := storage.NewSurveyQueryStub()

	_, err = storage.NewDB(DBPath)

	dictionariesStore := storage.NewDictionariesStore(draft_survey.Dictionaries)

	validator := validation.New()

	h := handler.New(userStore, surveyStore, surveyQueryStore, dictionariesStore, validator, version)

	handler.SetupRoutes(app, h)

	log.Fatal(app.Listen(":3399"))
}
