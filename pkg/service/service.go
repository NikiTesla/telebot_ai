package service

import (
	"fmt"
	"os"

	"telebotai/pkg/ai"
	"telebotai/pkg/telegram"

	log "github.com/sirupsen/logrus"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Run() error {
	openai_key := os.Getenv("OPENAI_API_KEY")
	if openai_key == "" {
		return fmt.Errorf("OPENAI_API_KEY env is empty")
	}
	aiService, err := ai.NewAiService(openai_key)
	if err != nil {
		return fmt.Errorf("failed to create AI Service, err: %w", err)
	}

	bot, err := telegram.NewBot(openai_key, aiService)
	if err != nil {
		return fmt.Errorf("failed connect telegram bot, err: %w", err)
	}
	return bot.Listen()
}

func (s *Service) Stop() {
	log.Debug("service was stopped")
}
