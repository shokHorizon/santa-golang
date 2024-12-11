package main

import (
	"fmt"
	"math/rand"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		panic(err)
	}

	userNames := make(map[int64]string)
	userIds := make(map[int64]int64)
	userWishes := make(map[int64]string)

	var owner int64

	bot.Debug = true

	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	updateConfig := tgbotapi.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := bot.GetUpdatesChan(updateConfig)

	// Let's go through each update that we're getting from Telegram.
	for update := range updates {
		// Telegram can send many types of updates depending on what your Bot
		// is up to. We only want to look at messages for now, so we can
		// discard any other updates.
		if update.Message == nil {
			continue
		}

		var msg tgbotapi.MessageConfig

		switch update.Message.Text {
		case "/own":
			owner = update.Message.From.ID
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Owner добавлен")
			userIds[owner] = 0
		case "/make":
			if owner == update.Message.From.ID {
				fillMap(userIds)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprint(userIds))
			}
		case "/random":
			if owner == update.Message.From.ID {
				newUser := rand.Int63()
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Добавлен user с id %d", newUser))
				userIds[newUser] = 0
			}
		case "/list":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprint(userIds))
		default:
			if id, ok := userIds[update.Message.From.ID]; !ok {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Вы добавлены в список участников, напишите одним сообщением свой вишлист")
				userIds[update.Message.From.ID] = 0
				userNames[update.Message.From.ID] = update.Message.From.UserName
			} else if id == 0 {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ваш новый вишлист\n%s\n\nВы можете изменить вишлист в любой момент, просто снова напишите сообщение", update.Message.Text))
				userWishes[update.Message.From.ID] = update.Message.Text
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ты даришь подарок: %s\nЕго вишлист: %s", userNames[userIds[update.Message.From.ID]], userWishes[userIds[update.Message.From.ID]]))
			}

		}

		// We'll also say that this message is a reply to the previous message.
		// For any other specifications than Chat ID or Text, you'll need to
		// set fields on the `MessageConfig`.

		// Okay, we're sending our message off! We don't care about the message
		// we just sent, so we'll discard it.
		if _, err := bot.Send(msg); err != nil {
			// Note that panics are a bad way to handle errors. Telegram can
			// have service outages or network errors, you should retry sending
			// messages or more gracefully handle failures.
			panic(err)
		}
	}

}

func fillMap(m map[int64]int64) {
	var firstElem, prevElem int64
	for k, _ := range m {
		if firstElem == 0 {
			firstElem = k
		}
		if prevElem != 0 {
			m[k] = prevElem
		}
		prevElem = k
	}

	m[firstElem] = prevElem
}
