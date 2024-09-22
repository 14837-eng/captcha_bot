package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken string `validate:"required"`
	EmojiCount       int    `validate:"required,min=1,max=10"`
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("ошибка при загрузке файла .env: %w", err)
	}

	config := &Config{
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		EmojiCount:       3, // значение по умолчанию
	}

	emojiCountStr := os.Getenv("EMOJI_COUNT")
	if emojiCountStr != "" {
		emojiCount, err := strconv.Atoi(emojiCountStr)
		if err != nil {
			return nil, fmt.Errorf("неверное значение EMOJI_COUNT: %w", err)
		}
		config.EmojiCount = emojiCount
	}

	validate := validator.New()
	err = validate.Struct(config)
	if err != nil {
		return nil, fmt.Errorf("ошибка валидации конфигурации: %w", err)
	}

	return config, nil
}
