package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/patrickmn/go-cache"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	// Start the telegramBot
	telegramBot, err := tgbotapi.NewBotAPI(os.Getenv("telegram.token"))
	if err != nil {
		log.Panic(err)
	}

	if debug, found := os.LookupEnv("telegram.debug"); found {
		telegramBot.Debug, _ = strconv.ParseBool(debug)
	}

	log.Printf("Authorized on account %s", telegramBot.Self.UserName)

	// Start the north sea surf bot
	token := initialiseMeetnet()

	bot := NorthSeaSurfBot{
		Config: Config{
			currentToken: token,
		},
		DataCache: DataCache{
			Remote: cache.New(5*time.Minute, 10*time.Minute),
		},
	}

	// Start the HTTP API
	http.HandleFunc("/", bot.Hello)
	http.HandleFunc("/health", bot.Hello)
	http.HandleFunc("/api/config", bot.Hello)
	http.HandleFunc("/api/meetnet/catalog", bot.Hello)
	http.HandleFunc("/api/cefas/current", bot.Current)
	http.HandleFunc("/api/cache", bot.Cache)

	go func() {
		log.Println("Starting http server on port 8080")
		http.ListenAndServe(":8080", nil)
	}()

	// Respond to commands
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

		for _, message := range bot.processCommand(&bot, update.Message) {
			if _, err := telegramBot.Send(message); err != nil {
				log.Panic(err)
			}
		}
	}
}

func initialiseMeetnet() Token {
	log.Println("Retrieving meetnet token")

	return login(os.Getenv("meetnet.user"), os.Getenv("meetnet.pass"))
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
