package main

import (
	"log"

	"github.com/14837-eng/captcha_bot/config"
	"github.com/14837-eng/captcha_bot/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	config.Init()

	bot, err := tgbotapi.NewBotAPI(config.Config.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = !config.Config.Production

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		telegram.HandleUpdate(bot, update)
	}
}
