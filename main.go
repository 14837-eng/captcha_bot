package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"captcha-bot/config"
	"captcha-bot/telegram"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %v", err)
	}

	// Создание нового бота
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Авторизован на аккаунте %s", bot.Self.UserName)

	// Настройка получения обновлений
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	// Получение канала обновлений
	updates := bot.GetUpdatesChan(updateConfig)

	fmt.Println("Бот успешно запущен!")

	// Создание обработчика событий
	handler := telegram.NewHandler(bot, cfg.EmojiCount)

	// Обработка обновлений
	for update := range updates {
		handler.HandleUpdate(update)
	}
}
