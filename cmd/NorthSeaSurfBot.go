package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
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

//

type NorthSeaSurfBot struct {
	Cache	BotCache
}

type BotCache struct {
	TokenCache		Token
	CatalogCache	Catalog
}

// Commands

func (t NorthSeaSurfBot) processCommand(input *tgbotapi.Message) []tgbotapi.MessageConfig {
	switch input.Command() {
	case "help":
		return []tgbotapi.MessageConfig{
			tgbotapi.NewMessage(input.Chat.ID, HELP),
		}
	case "locations":
		return []tgbotapi.MessageConfig{
			tgbotapi.NewMessage(input.Chat.ID, displayLocations(t.Cache.CatalogCache.Locations)),
		}
	case "availabledata":
		return []tgbotapi.MessageConfig{
			tgbotapi.NewMessage(input.Chat.ID, displayAvailableData(t.Cache.CatalogCache.AvailableData)),
		}
	case "currentdata":
		return []tgbotapi.MessageConfig{
			tgbotapi.NewMessage(input.Chat.ID, "test"),
		}
		//msg.Text = displayCurrentData(currentData(token))
	case "safekite":
		{
			id := "BL7WVC"
			data := currentDataForId(validateToken(t.Cache.TokenCache), []string{id})[id]
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
			current := getcurrent()
			id := input.CommandArguments()
			if feature, ok := current[id]; ok {
				return []tgbotapi.MessageConfig{
					tgbotapi.NewMessage(input.Chat.ID, feature.display()),
				}
			} else {
				return []tgbotapi.MessageConfig{
					tgbotapi.NewMessage(input.Chat.ID, "Could not find buoy: " + id),
				}
			}
		}
	case "cefasbuoys":
		return []tgbotapi.MessageConfig{
			tgbotapi.NewMessage(input.Chat.ID, displayBuoys(getcurrent())),
		}
	default:
		return []tgbotapi.MessageConfig{
			tgbotapi.NewMessage(input.Chat.ID, "I don't know that command"),
		}
	}

	return nil
}
