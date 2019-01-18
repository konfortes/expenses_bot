package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

// Bot ...
type Bot struct {
	*tgbotapi.BotAPI
}

func initBot() error {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		return errors.Wrap(err, "error creating bot instance")
	}

	// bot.Debug = true

	webhookURL, _ := url.Parse(fmt.Sprintf("%s/%s", os.Getenv("WEBHOOK_URL"), bot.Token))
	_, err = bot.SetWebhook(tgbotapi.WebhookConfig{URL: webhookURL})
	if err != nil {
		return errors.Wrap(err, "error setting web hook")
	}

	// info, err := bot.GetWebhookInfo()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if info.LastErrorDate != 0 {
	// 	log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	// }

	myBot := &Bot{bot}
	updates := bot.ListenForWebhook("/" + bot.Token)

	for update := range updates {
		myBot.handleUpdate(update)
	}

	return nil
}

func (bot Bot) handleUpdate(update tgbotapi.Update) {
	// TODO: this does not work. app crashes
	defer handlePanic()
	if update.Message == nil {
		return
	}

	jsonedUpdate, _ := json.Marshal(update)
	fmt.Println(string(jsonedUpdate))

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	switch update.Message.Command() {
	case "expense":
		command := strings.Split(update.Message.Text, " ")
		amount, _ := strconv.ParseFloat(command[1], 32)
		description := strings.Join(command[2:], " ")
		persistExpense(float32(amount), description, update.Message.From.ID)

		msg.Text = "how would you categorize this expense?"
		msg.ReplyMarkup = newCategoriesKeyboard()
	case "weather":
		msg.Text = weatherClient.current()
	default:
		msg.Text = "I only understand predefined commands. press the '/' button"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	}

	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}

func newCategoriesKeyboard() tgbotapi.ReplyKeyboardMarkup {
	var keyboard [][]tgbotapi.KeyboardButton

	keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Bills"),
		tgbotapi.NewKeyboardButton("Food"),
		tgbotapi.NewKeyboardButton("Transportation"),
	),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Pleasure"),
			tgbotapi.NewKeyboardButton("Vacation"),
			tgbotapi.NewKeyboardButton("Misc."),
		))

	return tgbotapi.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard:       keyboard,
	}
}
