package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/patrickmn/go-cache"
	"net/http"
)

// Constants

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

//

type NorthSeaSurfBot struct {
	Config    Config
	DataCache DataCache
}

type Config struct {
	currentToken Token
}

func (t Config) getToken() Token {
	return t.currentToken.validate()
}

type DataCache struct {
	Remote *cache.Cache
}

// HTTP API
func (t NorthSeaSurfBot) Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, I am the North Sea surf bot."))
}

func (t NorthSeaSurfBot) Current(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(getcurrent(&t))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (t NorthSeaSurfBot) ListCache(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(t.DataCache.Remote.Items())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (t NorthSeaSurfBot) FlushCache(w http.ResponseWriter, r *http.Request) {
	t.DataCache.Remote.Flush()
}

// Commands
func (t NorthSeaSurfBot) processCommand(bot *NorthSeaSurfBot, input *tgbotapi.Message) []tgbotapi.MessageConfig {
	switch input.Command() {
	case "help":
		return []tgbotapi.MessageConfig{
			tgbotapi.NewMessage(input.Chat.ID, HELP),
		}
	case "locations":
		return []tgbotapi.MessageConfig{
			tgbotapi.NewMessage(input.Chat.ID, displayLocations(catalog(bot).Locations)),
		}
	case "availabledata":
		return []tgbotapi.MessageConfig{
			tgbotapi.NewMessage(input.Chat.ID, displayAvailableData(catalog(bot).AvailableData)),
		}
	case "currentdata":
		return []tgbotapi.MessageConfig{
			tgbotapi.NewMessage(input.Chat.ID, "test"),
		}
		//msg.Text = displayCurrentData(currentData(token))
	case "safekite":
		{
			id := "BL7WVC"
			data := currentDataForId(bot, id)
			beaufort := MeterePerSecondToBeaufortScale(data.Value)
			text := "Wind: " + fmt.Sprintf("%.2f", data.Value) + "m/s (" + DisplayBeaufort(beaufort) + ")\n"
			if safeToKite(beaufort) {
				text += "It is safe to kite"
			} else {
				text += "It is not safe to kite"
			}

			return []tgbotapi.MessageConfig{
				tgbotapi.NewMessage(input.Chat.ID, text),
			}
		}
	case "cefas":
		{
			current := getcurrent(bot)
			id := input.CommandArguments()
			if feature, ok := current[id]; ok {
				return []tgbotapi.MessageConfig{
					tgbotapi.NewMessage(input.Chat.ID, feature.display()),
				}
			} else {
				return []tgbotapi.MessageConfig{
					tgbotapi.NewMessage(input.Chat.ID, "Could not find buoy: "+id),
				}
			}
		}
	case "cefasbuoys":
		return []tgbotapi.MessageConfig{
			tgbotapi.NewMessage(input.Chat.ID, displayBuoys(getcurrent(bot))),
		}
	default:
		return []tgbotapi.MessageConfig{
			tgbotapi.NewMessage(input.Chat.ID, "I don't know that command"),
		}
	}

	return nil
}
