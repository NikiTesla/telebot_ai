package service

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type TelegramBot struct {
	Bot *tgbotapi.BotAPI
	ai  *AiService
}

func NewTeleBot(token string, ai *AiService) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEBOT_API"))
	if err != nil {
		logrus.Fatalf("Cannot connect to telegram bot, err: %s", err.Error())
	}

	return &TelegramBot{bot, ai}, err
}

func (tB *TelegramBot) Listen() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates, err := tB.Bot.GetUpdatesChan(updateConfig)
	if err != nil {
		logrus.Printf("error occured while getting updates chain %s", err.Error())
	}

	for update := range updates {

		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		go func(tgbotapi.MessageConfig) {
			log.Println("start working")

			aiAnswer, err := tB.ai.MakeRequest(msg.Text)
			if err != nil {
				logrus.Printf("can't get answer from AI, err: %s", err.Error())
				return
			}

			msg.Text = aiAnswer

			if _, err := tB.Bot.Send(msg); err != nil {
				logrus.Printf("can't send message, err: %s", err.Error())
				return
			}

			log.Printf("result")
		}(msg)

	}
}
