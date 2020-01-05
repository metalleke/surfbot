package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"strconv"
)

func main() {
	// Start the webserver

	// Start the telegramBot
	telegramBot, err := tgbotapi.NewBotAPI(os.Getenv("telegram.token"))
	if err != nil {
		log.Panic(err)
	}

	if debug, found := os.LookupEnv("telegram.debug"); found {
		telegramBot.Debug, _ = strconv.ParseBool(debug)
	}

	log.Printf("Authorized on account %s", telegramBot.Self.UserName)

	token, catalog := initialiseMeetnet()
	northSeaSurfBot := NorthSeaSurfBot{
		Cache: BotCache{
			TokenCache:   token,
			CatalogCache: catalog,
		}}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := telegramBot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		for _, message := range northSeaSurfBot.processCommand(update.Message) {
			if _, err := telegramBot.Send(message); err != nil {
				log.Panic(err)
			}
		}
	}
}

func initialiseMeetnet() (Token, Catalog) {
	// Initialise meetnet data
	log.Println("Initializing meetnet data")

	token := login(os.Getenv("meetnet.user"), os.Getenv("meetnet.pass"))
	catalog := catalog(token)

	log.Println("Locations: ", catalog.Locations)
	log.Println("Parameters: ", catalog.Parameters)
	log.Println("ParameterTypes: ", catalog.ParameterTypes)
	log.Println("AvailableData: ", catalog.AvailableData)

	return token, catalog
}

func displayAvailableData(data map[string]AvailableData) string {
	result := ""
	for _, aData := range data {
		id := aData.Id
		location := ""
		parameter := ""

		result += id + " at " + location + " " + parameter + "\n"
	}
	return result
}
