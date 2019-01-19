package main

import (
	"encoding/json"
	"fmt"
	"log"
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

	bot.Debug = true

	webhookURL, _ := url.Parse(fmt.Sprintf("%s/%s", os.Getenv("WEBHOOK_URL"), bot.Token))
	fmt.Println("webhook url: ", webhookURL)
	_, err = bot.SetWebhook(tgbotapi.WebhookConfig{URL: webhookURL})
	if err != nil {
		return errors.Wrap(err, "error setting web hook")
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Print(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	myBot := &Bot{bot}
	updates := bot.ListenForWebhook("/" + bot.Token)

	for update := range updates {
		myBot.handleUpdate(update)
	}

	return nil
}

var ongoingExpenseReport = make(map[int]*Expense)

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
		expense := Expense{}
		command := strings.Split(update.Message.Text, " ")
		amount, _ := strconv.ParseFloat(command[1], 32)
		expense.Amount = float32(amount)
		expense.Description = strings.Join(command[2:], " ")
		expense.UserID = update.Message.From.ID
		ongoingExpenseReport[expense.UserID] = &expense

		msg.Text = "how would you categorize this expense?"
		msg.ReplyMarkup = newCategoriesKeyboard()
	case "weather":
		msg.Text = weatherClient.current()
	default:
		expense := ongoingExpenseReport[update.Message.From.ID]
		if expense != nil {
			if validCategory(update.Message.Text) {
				expense.Category = update.Message.Text
				if err := persistExpense(expense); err != nil {
					log.Println("error persisting expense: ", err)
				}
				msg.Text = "expense persisted successfully"
			} else {
				msg.Text = "unrecognized category. expense could not be persisted"
			}
			delete(ongoingExpenseReport, update.Message.From.ID)
		}
	}

	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}

func validCategory(category string) bool {
	for _, c := range []string{"Bills", "Food", "Transportation", "Pleasure", "Vacation", "Misc."} {
		if c == category {
			return true
		}
	}
	return false
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
		ResizeKeyboard:  true,
		Keyboard:        keyboard,
		OneTimeKeyboard: true,
	}
}
