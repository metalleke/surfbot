package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
)

const HELP = `
help: this message
locations: list of all locations
availabledata: list of all data points
currentdata <ID>: the current data of a buoy
getdata <ID> <FROM> <TO>: returns the data of a specified buoy in a time range
safekite: current wind value for safe kiting
cefas <ID>: current cefas data
cefasbuoys: list of all cefas buoys
`

func main() {
	// Start the bot
	bot, err := tgbotapi.NewBotAPI(os.Getenv("telegram.token"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	token, catalog := initialiseMeetnet()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = HELP
		case "locations":
			msg.Text = displayLocations(catalog.Locations)
		case "availabledata":
			msg.Text = displayAvailableData(catalog.AvailableData)
		case "currentdata":
			msg.Text = "test"
			//msg.Text = displayCurrentData(currentData(token))
		case "safekite": {
			data := currentDataForId(validateToken(token), []string{"BL7WVC"})
			msg.Text = displayCurrentData(catalog, data) + " (" + MeterePerSecondToBeaufortScale(data["BL7WVC"].Value) + ")"
		}
		case "cefas": {
			current := getcurrent()
			id := update.Message.CommandArguments()
			if feature, ok := current[id]; ok {
				msg.Text = displayCurrentWaveHeight(feature)
			} else  {
				msg.Text = "Could not find buoy: " + id
			}
		}
		case "cefasbuoys":
			msg.Text = displayBuoys(getcurrent())
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
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
