package telegram

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/14837-eng/captcha_bot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	defaultEmojiList = []string{"😀", "😃", "😄", "😁", "😆", "😅", "😂", "🤣", "😊", "😇", "🙂", "🙃", "😉", "😌", "😍", "🥰", "😘", "😗", "😙", "😚", "😋", "😛", "😝", "😜", "🤪", "🤨", "🧐", "🤓", "😎", "🤩", "🥳"}
	userCaptchas     = make(map[int64]captchaInfo)
)

type captchaInfo struct {
	captcha   string
	startTime time.Time
	messageID int
}

func getEmojiList() []string {
	if len(config.Config.CustomEmojiList) > 0 {
		return config.Config.CustomEmojiList
	}
	return defaultEmojiList
}

func generateCaptcha(count int) string {
	rand.Seed(time.Now().UnixNano())
	emojiList := getEmojiList()
	captcha := ""
	for i := 0; i < count; i++ {
		captcha += emojiList[rand.Intn(len(emojiList))]
	}
	return captcha
}

func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message != nil {
		// Обработка новых участников
		if update.Message.NewChatMembers != nil {
			for _, newUser := range update.Message.NewChatMembers {
				// Проверяем, что новый пользователь не является самим ботом
				if newUser.ID == bot.Self.ID {
					continue // Пропускаем обработку, если это бот
				}

				captcha := generateCaptcha(config.Config.EmojiCount)
				userCaptchas[newUser.ID] = captchaInfo{
					captcha:   captcha,
					startTime: time.Now(),
				}

				// Создаем клавиатуру с эмодзи
				var keyboard [][]tgbotapi.InlineKeyboardButton
				buttonCount := config.Config.EmojiButtonCount
				if buttonCount < config.Config.EmojiCount {
					buttonCount = config.Config.EmojiCount
				}
				emojiList := getEmojiList()
				if buttonCount > len(emojiList) {
					buttonCount = len(emojiList)
				}

				// Создаем список эмодзи для кнопок, начиная с эмодзи из капчи
				buttonEmojis := []string{}
				for _, emoji := range captcha {
					buttonEmojis = append(buttonEmojis, string(emoji))
				}

				// Добавляем случайные эмодзи, пока не достигнем нужного количества кнопок
				remainingEmojis := make([]string, len(emojiList))
				copy(remainingEmojis, emojiList)
				rand.Shuffle(len(remainingEmojis), func(i, j int) {
					remainingEmojis[i], remainingEmojis[j] = remainingEmojis[j], remainingEmojis[i]
				})

				for _, emoji := range remainingEmojis {
					if len(buttonEmojis) >= buttonCount {
						break
					}
					if !strings.Contains(captcha, emoji) {
						buttonEmojis = append(buttonEmojis, emoji)
					}
				}

				// Перемешиваем список эмодзи для кнопок
				rand.Shuffle(len(buttonEmojis), func(i, j int) {
					buttonEmojis[i], buttonEmojis[j] = buttonEmojis[j], buttonEmojis[i]
				})

				// Располагаем кнопки горизонтально, максимум 5 кнопок в ряду
				buttonsPerRow := 5
				var row []tgbotapi.InlineKeyboardButton
				for i, emoji := range buttonEmojis {
					button := tgbotapi.NewInlineKeyboardButtonData(emoji, fmt.Sprintf("captcha:%d:%s", newUser.ID, emoji))
					row = append(row, button)

					if len(row) == buttonsPerRow || i == len(buttonEmojis)-1 {
						keyboard = append(keyboard, row)
						row = []tgbotapi.InlineKeyboardButton{}
					}
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(config.Config.WelcomeMessage, newUser.FirstName, config.Config.CaptchaTimeoutMinutes/time.Minute, captcha))
				msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboard}
				sentMsg, err := bot.Send(msg)
				if err != nil {
					log.Printf("Failed to send captcha message: %v", err)
					return
				}

				userCaptchas[newUser.ID] = captchaInfo{
					captcha:   captcha,
					startTime: time.Now(),
					messageID: sentMsg.MessageID,
				}

				// Ограничиваем права пользователя
				restrictChatMember := tgbotapi.RestrictChatMemberConfig{
					ChatMemberConfig: tgbotapi.ChatMemberConfig{
						ChatID: update.Message.Chat.ID,
						UserID: newUser.ID,
					},
					Permissions: &tgbotapi.ChatPermissions{
						CanSendMessages: false,
					},
				}
				bot.Request(restrictChatMember)

				// Запускаем горутину для проверки таймаута
				go checkCaptchaTimeout(bot, update.Message.Chat.ID, newUser.ID, sentMsg.MessageID)
			}
		}
	} else if update.CallbackQuery != nil {
		// Обработка нажатий на кнопки капчи
		data := strings.Split(update.CallbackQuery.Data, ":")
		if len(data) == 3 && data[0] == "captcha" {
			userID := update.CallbackQuery.From.ID
			clickedEmoji := data[2]

			if currentCaptcha, ok := userCaptchas[userID]; ok {
				if strings.HasPrefix(currentCaptcha.captcha, clickedEmoji) {
					userCaptchas[userID] = captchaInfo{
						captcha:   currentCaptcha.captcha[len(clickedEmoji):],
						startTime: currentCaptcha.startTime,
						messageID: currentCaptcha.messageID,
					}
					if userCaptchas[userID].captcha == "" {
						// Капча пройдена успешно
						delete(userCaptchas, userID)

						// Возвращаем полные права пользователю
						unrestrictChatMember := tgbotapi.RestrictChatMemberConfig{
							ChatMemberConfig: tgbotapi.ChatMemberConfig{
								ChatID: update.CallbackQuery.Message.Chat.ID,
								UserID: userID,
							},
							Permissions: &tgbotapi.ChatPermissions{
								CanSendMessages:       true,
								CanSendMediaMessages:  true,
								CanSendOtherMessages:  true,
								CanAddWebPagePreviews: true,
							},
						}
						bot.Request(unrestrictChatMember)

						// Удаляем сообщение с кнопками
						deleteMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
						bot.Request(deleteMsg)

						callbackConfig := tgbotapi.CallbackConfig{
							CallbackQueryID: update.CallbackQuery.ID,
							Text:            config.Config.CaptchaSuccessMessage,
							ShowAlert:       true,
						}
						bot.Request(callbackConfig)

						bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf(config.Config.CaptchaPassedAnnouncement, update.CallbackQuery.From.UserName)))
					} else {
						callbackConfig := tgbotapi.CallbackConfig{
							CallbackQueryID: update.CallbackQuery.ID,
							Text:            config.Config.CaptchaPartialSuccessMessage,
							ShowAlert:       false,
						}
						bot.Request(callbackConfig)
					}
				} else {
					callbackConfig := tgbotapi.CallbackConfig{
						CallbackQueryID: update.CallbackQuery.ID,
						Text:            config.Config.CaptchaFailMessage,
						ShowAlert:       true,
					}
					bot.Request(callbackConfig)
				}
			}
		}
	}
}

func checkCaptchaTimeout(bot *tgbotapi.BotAPI, chatID int64, userID int64, messageID int) {
	time.Sleep(config.Config.CaptchaTimeoutMinutes)

	if _, ok := userCaptchas[userID]; ok {
		// Если капча все еще не пройдена, кикаем пользователя
		delete(userCaptchas, userID)

		// Удаляем сообщение с кнопками капчи
		deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
		_, err := bot.Request(deleteMsg)
		if err != nil {
			log.Printf("Failed to delete captcha message: %v", err)
		}

		kickChatMember := tgbotapi.KickChatMemberConfig{
			ChatMemberConfig: tgbotapi.ChatMemberConfig{
				ChatID: chatID,
				UserID: userID,
			},
		}

		_, err = bot.Request(kickChatMember)
		if err != nil {
			// Если не удалось кикнуть пользователя, отправляем сообщение
			bot.Send(tgbotapi.NewMessage(chatID, config.Config.KickFailMessage))
		} else {
			bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf(config.Config.KickSuccessMessage, userID)))
		}
	}
}
