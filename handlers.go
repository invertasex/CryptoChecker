package main

import (
	"cryptobot/utils"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	switch {
	case msg.Text == "/start":
		reply := "👋 Hello! I'm CryptoBot. Here's what I can do:\n\n/price BTC — check BTC price\n/balance [wallet] — check USDT balance\n/help — command list"
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, reply))

	case msg.Text == "/help":
		reply := "/price [TICKER] — get crypto price\n/balance [WALLET] — check USDT balance on TRON (TRC-20)"
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, reply))

	case strings.HasPrefix(msg.Text, "/price"):
		args := strings.Split(msg.Text, " ")
		if len(args) < 2 {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Please specify a ticker, e.g., /price BTC"))
			return
		}
		price, err := utils.GetPrice(strings.ToUpper(args[1]))
		if err != nil {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Failed to retrieve price."))
			return
		}
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Price of %s: $%.2f", args[1], price)))

	case strings.HasPrefix(msg.Text, "/balance"):
		args := strings.Split(msg.Text, " ")
		if len(args) < 2 {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Please specify a wallet address, e.g., /balance TH..."))
			return
		}
		balance, err := utils.GetTRC20Balance(args[1])
		if err != nil {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Failed to retrieve balance."))
			return
		}
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("USDT (TRC20) Balance: %.2f USDT", balance)))

	default:
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Unknown command. Type /help"))
	}
}
