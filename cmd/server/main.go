package main

import (
	"log"

	"github.com/AVZotov/draft-survey/internal/handler"
	"github.com/AVZotov/draft-survey/internal/storage"
	"github.com/gofiber/fiber/v2"
)

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

	dictionariesStore, err := storage.NewDictionariesStore("./dictionaries")
	if err != nil {
		log.Fatal(err)
	}

	h := handler.New(userStore, surveyStore, dictionariesStore)

	handler.SetupRoutes(app, h)

	log.Fatal(app.Listen(":3399"))
}
