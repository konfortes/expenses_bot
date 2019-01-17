package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"

	// auto load .env file
	_ "github.com/joho/godotenv/autoload"
)

// Bot ...
type Bot struct {
	*tgbotapi.BotAPI
}

func initBot() (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		return nil, errors.Wrap(err, "error creating bot instance")
	}

	// bot.Debug = true

	webhookURL, _ := url.Parse(fmt.Sprintf("%s/%s", os.Getenv("WEBHOOK_URL"), bot.Token))
	_, err = bot.SetWebhook(tgbotapi.WebhookConfig{URL: webhookURL})
	if err != nil {
		return nil, errors.Wrap(err, "error setting web hook")
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	return &Bot{bot}, nil
}

var weatherClient WeatherClient

func main() {
	weatherClient = WeatherClient{URL: os.Getenv("WEATHER_API_URL"), Key: os.Getenv("WEATHER_API_KEY")}
	bot, err := initBot()
	if err != nil {
		log.Fatal(err)
		return
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe("localhost:8080", nil)

	for update := range updates {
		bot.handleUpdate(update)
	}
}

func (bot Bot) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	switch update.Message.Command() {
	case "expense":
		msg.Text = "how would you categorize this expense?"
		msg.ReplyMarkup = newCategoriesKeyboard()
	case "weather":
		msg.Text = weatherClient.current()
	case "test":

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
