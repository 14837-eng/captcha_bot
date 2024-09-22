package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Configuration struct {
	TelegramBotToken             string
	EmojiCount                   int
	EmojiButtonCount             int
	CaptchaTimeoutMinutes        time.Duration
	WelcomeMessage               string
	CaptchaSuccessMessage        string
	CaptchaPartialSuccessMessage string
	CaptchaFailMessage           string
	CaptchaPassedAnnouncement    string
	KickFailMessage              string
	KickSuccessMessage           string
	CustomEmojiList              []string
	Production                   bool
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

	captchaTimeoutMinutes, err := strconv.Atoi(os.Getenv("CAPTCHA_TIMEOUT_MINUTES"))
	if err != nil {
		log.Fatal("CAPTCHA_TIMEOUT_MINUTES must be a valid integer in .env file")
	}
	Config.CaptchaTimeoutMinutes = time.Duration(captchaTimeoutMinutes) * time.Minute

	Config.WelcomeMessage = os.Getenv("WELCOME_MESSAGE")
	Config.CaptchaSuccessMessage = os.Getenv("CAPTCHA_SUCCESS_MESSAGE")
	Config.CaptchaPartialSuccessMessage = os.Getenv("CAPTCHA_PARTIAL_SUCCESS_MESSAGE")
	Config.CaptchaFailMessage = os.Getenv("CAPTCHA_FAIL_MESSAGE")
	Config.CaptchaPassedAnnouncement = os.Getenv("CAPTCHA_PASSED_ANNOUNCEMENT")
	Config.KickFailMessage = os.Getenv("KICK_FAIL_MESSAGE")
	Config.KickSuccessMessage = os.Getenv("KICK_SUCCESS_MESSAGE")

	customEmojiList := os.Getenv("CUSTOM_EMOJI_LIST")
	if customEmojiList != "" {
		Config.CustomEmojiList = strings.Split(customEmojiList, ",")
	}

	production, err := strconv.ParseBool(os.Getenv("PRODUCTION"))
	if err != nil {
		log.Fatal("PRODUCTION must be a valid boolean (true/false) in .env file")
	}
	Config.Production = production
}
