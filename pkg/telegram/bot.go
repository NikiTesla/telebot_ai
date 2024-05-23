package telegram

import (
	"fmt"
	"os"
	"telebotai/pkg/ai"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type Bot struct {
	bot       *tgbotapi.BotAPI
	aiService *ai.Service
}

func NewBot(token string, aiService *ai.Service) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEBOT_API"))
	if err != nil {
		log.WithError(err).Fatal("failed to connect to telegram bot")
	}
	return &Bot{
		bot:       bot,
		aiService: aiService,
	}, err
}

func (b *Bot) Listen() error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates, err := b.bot.GetUpdatesChan(updateConfig)
	if err != nil {
		return fmt.Errorf("filed to get updates, err: %w", err)
	}

	for update := range updates {
		if update.Message == nil {
			log.Debug("update message is nil")
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		go func() {
			if err := b.processMessage(msg); err != nil {
				log.WithError(err).Error("failed to process message")
			}
		}()
	}
	return nil
}

func (b *Bot) processMessage(msg tgbotapi.MessageConfig) error {
	log.Debugf("processing message: %s", msg.Text)
	aiAnswer, err := b.aiService.MakeRequest(msg.Text)
	if err != nil {
		return fmt.Errorf("failed to make request, err: %w", err)
	}
	msg.Text = aiAnswer
	if _, err := b.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send message, err: %w", err)
	}
	return nil
}
