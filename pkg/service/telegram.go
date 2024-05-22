package service

import (
	"fmt"
	"os"
	"telebotai/pkg/ai"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type TelegramBot struct {
	bot       *tgbotapi.BotAPI
	aiService *ai.Service
}

func NewTelegramBot(token string, aiService *ai.Service) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEBOT_API"))
	if err != nil {
		log.WithError(err).Fatal("failed to connect to telegram bot")
	}
	return &TelegramBot{
		bot:       bot,
		aiService: aiService,
	}, err
}

func (tB *TelegramBot) Listen() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates, err := tB.bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.WithError(err).Fatal("filed to get updates")
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		go tB.processMessage(msg)
	}
}

func (tB *TelegramBot) processMessage(msg tgbotapi.MessageConfig) error {
	aiAnswer, err := tB.aiService.MakeRequest(msg.Text)
	if err != nil {
		return fmt.Errorf("failed to make request, err: %w", err)
	}
	msg.Text = aiAnswer
	if _, err := tB.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send message, err: %w", err)
	}
	return nil
}
