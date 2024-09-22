package telegram

import (
	"fmt"
	"log"
	"math/rand"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot *tgbotapi.BotAPI
}

func NewHandler(bot *tgbotapi.BotAPI) *Handler {
	return &Handler{bot: bot}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	log.Printf("%s", update.Message)

	if update.Message.NewChatMembers != nil {
		h.handleNewChatMembers(update.Message)
	}
}

func (h *Handler) handleNewChatMembers(message *tgbotapi.Message) {
	for _, newMember := range message.NewChatMembers {
		if newMember.ID == h.bot.Self.ID {
			continue
		}

		h.restrictNewMember(message.Chat.ID, &newMember)
		h.sendCaptchaMessage(message.Chat.ID, &newMember)
	}
}

func (h *Handler) restrictNewMember(chatID int64, member *tgbotapi.User) {
	log.Printf("Новый участник присоединился: %s", member.UserName)

	chatPermissions := tgbotapi.ChatPermissions{
		CanSendMessages:       false,
		CanSendMediaMessages:  false,
		CanSendOtherMessages:  false,
		CanAddWebPagePreviews: false,
	}

	restrictConfig := tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: member.ID,
		},
		Permissions: &chatPermissions,
	}

	_, err := h.bot.Request(restrictConfig)
	if err != nil {
		log.Printf("Ошибка при ограничении прав пользователя: %v", err)
	} else {
		log.Printf("Пользователь %s успешно ограничен в правах", member.UserName)
	}
}

func (h *Handler) sendCaptchaMessage(chatID int64, member *tgbotapi.User) {
	keyboard := h.createEmojiKeyboard()
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("@%s, вы были лишены права голоса. Чтобы продолжить, вам нужно пройти капчу.", member.UserName))
	msg.ReplyMarkup = keyboard
	h.bot.Send(msg)
}

func (h *Handler) createEmojiKeyboard() tgbotapi.InlineKeyboardMarkup {
	emojis := []string{"😀", "😎", "🤔", "🎉", "🌈", "🍕", "🐱", "🚀", "🌺"}
	rand.Shuffle(len(emojis), func(i, j int) { emojis[i], emojis[j] = emojis[j], emojis[i] })

	var keyboard [][]tgbotapi.InlineKeyboardButton
	row := []tgbotapi.InlineKeyboardButton{}
	for i := 0; i < 3; i++ {
		button := tgbotapi.NewInlineKeyboardButtonData(emojis[i], fmt.Sprintf("captcha_%d", i))
		row = append(row, button)
	}
	keyboard = append(keyboard, row)

	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}
