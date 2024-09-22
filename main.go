package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"captcha-bot/telegram"

	"github.com/joho/godotenv"
)

func main() {
	// Загрузка переменных окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка при загрузке файла .env")
	}

	// Получение токена бота из переменных окружения
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Токен бота не найден в переменных окружения")
	}

	// Создание нового бота
	bot, err := tgbotapi.NewBotAPI(botToken)
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
	handler := telegram.NewHandler(bot)

	// Обработка обновлений
	for update := range updates {
		handler.HandleUpdate(update)
	}
}
