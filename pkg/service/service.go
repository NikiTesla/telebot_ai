package service

import (
	"os"

	"telebotai/pkg/ai"

	log "github.com/sirupsen/logrus"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) MustRun() {
	openai_key := os.Getenv("OPENAI_API_KEY")
	if openai_key == "" {
		log.Fatal("OPENAI_API_KEY env is empty")
	}
	aiService, err := ai.NewAiService(openai_key, "configs/api_request.json")
	if err != nil {
		log.WithError(err).Fatalf("failed to create AI Service")
	}

	bot, err := NewTelegramBot(openai_key, aiService)
	if err != nil {
		log.WithError(err).Fatalf("failed connect telegram bot")
	}
	bot.Listen()
}

func (s *Service) Stop() {
	log.Debug("service was stopped")
}
