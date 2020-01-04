package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func processCommand(input tgbotapi.Message, response tgbotapi.MessageConfig) {
	// Extract the command from the Message.
	switch input.Command() {
	case "help":
		response.Text = "type /sayhi or /status."
	case "catalog":
		response.Text = "Hi :)"
	case "buoy":
		response.Text = "I'm ok."
	default:
		response.Text = "I don't know that command"
	}
}
