package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
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
		return nil, err
	}

	bot.Debug = true

	webhookURL, _ := url.Parse(fmt.Sprintf("%s/%s", os.Getenv("WEBHOOK_URL"), bot.Token))
	_, err = bot.SetWebhook(tgbotapi.WebhookConfig{URL: webhookURL})
	if err != nil {
		return nil, err
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

func main() {
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
		msg.Text = "what did you spent your money on?"
	case "weather":
		msg.Text = "It is 14 celsius and 60% raining"
	default:
		msg.Text = "I only understand predefined commands. press the '/' button"
	}
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}
