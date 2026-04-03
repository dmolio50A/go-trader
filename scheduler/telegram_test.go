package main

import (
	"strings"
	"testing"
)

func TestFormatTradeDMPlain_NoMarkdown(t *testing.T) {
	sc := StrategyConfig{ID: "hl-sma-btc", Platform: "hyperliquid", Type: "perps"}
	trade := Trade{
		Symbol:   "BTC",
		Side:     "buy",
		Quantity: 0.15,
		Price:    67845.00,
		Value:    10176.75,
		Details:  "Open long 0.150000 @ $67845.00 (fee $10.18)",
	}
	msg := FormatTradeDMPlain(sc, trade, "paper")

	if strings.Contains(msg, "**") {
		t.Errorf("plain format should not contain ** markdown, got:\n%s", msg)
	}
	if !strings.Contains(msg, "TRADE EXECUTED") {
		t.Errorf("expected 'TRADE EXECUTED', got:\n%s", msg)
	}
	if !strings.Contains(msg, "hl-sma-btc") {
		t.Errorf("expected strategy ID, got:\n%s", msg)
	}
	if !strings.Contains(msg, "BUY") {
		t.Errorf("expected BUY, got:\n%s", msg)
	}
	if !strings.Contains(msg, "Mode: paper") {
		t.Errorf("expected 'Mode: paper', got:\n%s", msg)
	}
}

func TestFormatTradeDMPlain_CloseTrade(t *testing.T) {
	sc := StrategyConfig{ID: "hl-rmc-eth", Platform: "hyperliquid", Type: "perps"}
	trade := Trade{
		Symbol:   "ETH",
		Side:     "sell",
		Quantity: 0.47,
		Price:    3077.70,
		Value:    1446.52,
		Details:  "Close long, PnL: $34.35 (fee $1.23)",
	}
	msg := FormatTradeDMPlain(sc, trade, "live")

	if strings.Contains(msg, "**") {
		t.Errorf("plain format should not contain ** markdown, got:\n%s", msg)
	}
	if !strings.Contains(msg, "TRADE CLOSED") {
		t.Errorf("expected 'TRADE CLOSED', got:\n%s", msg)
	}
	if !strings.Contains(msg, "PnL: $34.35") {
		t.Errorf("expected PnL in close trade, got:\n%s", msg)
	}
	if !strings.Contains(msg, "Mode: live") {
		t.Errorf("expected 'Mode: live', got:\n%s", msg)
	}
}

func TestFormatTradeDMPlain_VsDiscord(t *testing.T) {
	sc := StrategyConfig{ID: "hl-sma-btc", Platform: "hyperliquid", Type: "perps"}
	trade := Trade{
		Symbol:   "BTC",
		Side:     "buy",
		Quantity: 0.15,
		Price:    67845.00,
		Value:    10176.75,
		Details:  "Open long 0.150000 @ $67845.00 (fee $10.18)",
	}

	discord := FormatTradeDM(sc, trade, "paper")
	telegram := FormatTradeDMPlain(sc, trade, "paper")

	// Discord version has ** bold markers
	if !strings.Contains(discord, "**") {
		t.Errorf("Discord format should contain ** markdown")
	}
	// Telegram version has no ** markers
	if strings.Contains(telegram, "**") {
		t.Errorf("Telegram format should not contain ** markdown")
	}
	// Both contain the same core data
	for _, want := range []string{"TRADE EXECUTED", "hl-sma-btc", "BUY", "Mode: paper"} {
		if !strings.Contains(discord, want) {
			t.Errorf("Discord format missing %q", want)
		}
		if !strings.Contains(telegram, want) {
			t.Errorf("Telegram format missing %q", want)
		}
	}
}
