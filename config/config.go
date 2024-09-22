package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Configuration struct {
	TelegramBotToken string
	EmojiCount       int
	EmojiButtonCount int
}

var Config Configuration

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Config.TelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	if Config.TelegramBotToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set in .env file")
	}

	emojiCount, err := strconv.Atoi(os.Getenv("EMOJI_COUNT"))
	if err != nil {
		log.Fatal("EMOJI_COUNT must be a valid integer in .env file")
	}
	Config.EmojiCount = emojiCount

	emojiButtonCount, err := strconv.Atoi(os.Getenv("EMOJI_BUTTON_COUNT"))
	if err != nil {
		log.Fatal("EMOJI_BUTTON_COUNT must be a valid integer in .env file")
	}
	Config.EmojiButtonCount = emojiButtonCount
}
