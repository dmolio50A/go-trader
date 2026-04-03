package main

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramNotifier wraps a Telegram Bot API client for sending messages.
type TelegramNotifier struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

// NewTelegramNotifier creates a Telegram bot client and verifies the token.
func NewTelegramNotifier(token string, chatID int64) (*TelegramNotifier, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("create telegram bot: %w", err)
	}
	return &TelegramNotifier{bot: bot, chatID: chatID}, nil
}

// SendMessage sends a text message to the configured chat ID.
func (t *TelegramNotifier) SendMessage(content string) error {
	msg := tgbotapi.NewMessage(t.chatID, content)
	_, err := t.bot.Send(msg)
	return err
}

// FormatTradeDMPlain formats a Trade into a plain-text DM (no Discord markdown).
func FormatTradeDMPlain(sc StrategyConfig, trade Trade, mode string) string {
	isClose := strings.Contains(trade.Details, "Close")

	icon := "🟢"
	header := "TRADE EXECUTED"
	if isClose {
		icon = "🔴"
		header = "TRADE CLOSED"
	}

	platformLabel := strings.ToUpper(sc.Platform[:1]) + sc.Platform[1:]
	typeLabel := sc.Type

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s %s\n", icon, header))
	sb.WriteString(fmt.Sprintf("Strategy: %s (%s %s)\n", sc.ID, platformLabel, typeLabel))
	sb.WriteString(fmt.Sprintf("%s — %s %.6g @ $%s\n", trade.Symbol, strings.ToUpper(trade.Side), trade.Quantity, fmtComma(trade.Price)))

	valueLine := fmt.Sprintf("Value: $%s", fmtComma(trade.Value))
	if isClose {
		if idx := strings.Index(trade.Details, "PnL: $"); idx >= 0 {
			pnlStr := trade.Details[idx+len("PnL: $"):]
			if end := strings.Index(pnlStr, " "); end >= 0 {
				pnlStr = pnlStr[:end]
			}
			valueLine += fmt.Sprintf(" | PnL: $%s", pnlStr)
		}
	}
	valueLine += fmt.Sprintf(" | Mode: %s", mode)
	sb.WriteString(valueLine)

	return sb.String()
}
