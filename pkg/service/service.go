package service

import (
	"os"

	"github.com/sirupsen/logrus"
)

type AI interface {
	MakeRequest(text string) (answer string, err error)
}

type TeleBot interface {
	Listen()
}

type Service struct {
	AI
	TeleBot
}

func (s *Service) Run() {
	aiService, err := CreateAiService(os.Getenv("OPENAI_API_KEY"), "configs/api_request.json")
	if err != nil {
		logrus.Fatalf("Cannot create AI Service %s", err.Error())
	}

	bot, err := NewTeleBot(os.Getenv("OPENAI_API_KEY"), aiService)
	if err != nil {
		logrus.Fatalf("Can't connect telegram bot: %s", err.Error())
	}

	bot.Listen()
}

func (s *Service) Stop() {

}
