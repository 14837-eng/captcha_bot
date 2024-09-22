package telegram

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/your-username/your-project/config"
)

var (
	emojiList    = []string{"😀", "😃", "😄", "😁", "😆", "😅", "😂", "🤣", "😊", "😇", "🙂", "🙃", "😉", "😌", "😍", "🥰", "😘", "😗", "😙", "😚", "😋", "😛", "😝", "😜", "🤪", "🤨", "🧐", "🤓", "😎", "🤩", "🥳"}
	userCaptchas = make(map[int64]string)
)

func generateCaptcha(count int) string {
	rand.Seed(time.Now().UnixNano())
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
				captcha := generateCaptcha(config.Config.EmojiCount)
				userCaptchas[newUser.ID] = captcha

				// Создаем клавиатуру с эмодзи
				var keyboard [][]tgbotapi.InlineKeyboardButton
				buttonCount := config.Config.EmojiButtonCount
				if buttonCount > len(emojiList) {
					buttonCount = len(emojiList)
				}
				shuffledEmojis := make([]string, len(emojiList))
				copy(shuffledEmojis, emojiList)
				rand.Shuffle(len(shuffledEmojis), func(i, j int) {
					shuffledEmojis[i], shuffledEmojis[j] = shuffledEmojis[j], shuffledEmojis[i]
				})
				for i := 0; i < buttonCount; i++ {
					emoji := shuffledEmojis[i]
					keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{
						tgbotapi.NewInlineKeyboardButtonData(emoji, fmt.Sprintf("captcha:%d:%s", newUser.ID, emoji)),
					})
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Добро пожаловать, %s! Пожалуйста, введите следующую капчу, нажимая на кнопки в правильном порядке:\n%s", newUser.FirstName, captcha))
				msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboard}
				bot.Send(msg)

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
			}
		}
	} else if update.CallbackQuery != nil {
		// Обработка нажатий на кнопки капчи
		data := strings.Split(update.CallbackQuery.Data, ":")
		if len(data) == 3 && data[0] == "captcha" {
			userID := update.CallbackQuery.From.ID
			clickedEmoji := data[2]

			if captcha, ok := userCaptchas[userID]; ok {
				if strings.HasPrefix(captcha, clickedEmoji) {
					userCaptchas[userID] = captcha[len(clickedEmoji):]
					if userCaptchas[userID] == "" {
						// Капча пройдена успешно
						delete(userCaptchas, userID)

						// Возвращаем права пользователю
						unrestrictChatMember := tgbotapi.RestrictChatMemberConfig{
							ChatMemberConfig: tgbotapi.ChatMemberConfig{
								ChatID: update.CallbackQuery.Message.Chat.ID,
								UserID: userID,
							},
							Permissions: &tgbotapi.ChatPermissions{
								CanSendMessages: true,
							},
						}
						bot.Request(unrestrictChatMember)

						// Используем bot.Request вместо несуществующего метода AnswerCallbackQuery
						callbackConfig := tgbotapi.CallbackConfig{
							CallbackQueryID: update.CallbackQuery.ID,
							Text:            "Капча пройдена успешно! Теперь вы можете писать в чате.",
							ShowAlert:       true,
						}
						bot.Request(callbackConfig)

						bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Пользователь @%s успешно прошел капчу!", update.CallbackQuery.From.UserName)))
					} else {
						callbackConfig := tgbotapi.CallbackConfig{
							CallbackQueryID: update.CallbackQuery.ID,
							Text:            "Правильно! Продолжайте.",
							ShowAlert:       false,
						}
						bot.Request(callbackConfig)
					}
				} else {
					callbackConfig := tgbotapi.CallbackConfig{
						CallbackQueryID: update.CallbackQuery.ID,
						Text:            "Неправильный порядок. Попробуйте еще раз.",
						ShowAlert:       true,
					}
					bot.Request(callbackConfig)
				}
			}
		}
	}
}
