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
	// Initialisation
	telegram := initialiseTelegram()
	token := initialiseMeetnet()
	bot := initialiseBot(token)
	initialiseHttpServer(bot)

	// Respond to commands
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := telegram.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		for _, message := range bot.processCommand(bot, update.Message) {
			if _, err := telegram.Send(message); err != nil {
				log.Panic(err)
			}
		}
	}
}

func initialiseTelegram() *tgbotapi.BotAPI {
	telegramBot, err := tgbotapi.NewBotAPI(os.Getenv("telegram.token"))
	if err != nil {
		log.Panic(err)
	}

	if debug, found := os.LookupEnv("telegram.debug"); found {
		telegramBot.Debug, _ = strconv.ParseBool(debug)
	}

	log.Printf("Telegram authorized on account %s", telegramBot.Self.UserName)

	return telegramBot
}

func initialiseMeetnet() Token {
	log.Printf("Retrieving meetnet token for %s", os.Getenv("meetnet.user"))

	return login(os.Getenv("meetnet.user"), os.Getenv("meetnet.pass"))
}

func initialiseBot(token Token) *NorthSeaSurfBot {
	return &NorthSeaSurfBot{
		Config: Config{
			currentToken: token,
		},
		DataCache: DataCache{
			Remote: cache.New(5*time.Minute, 10*time.Minute),
		},
	}
}

func initialiseHttpServer(bot *NorthSeaSurfBot) {
	http.HandleFunc("/", bot.Hello)

	http.HandleFunc("/health", bot.Health)

	http.HandleFunc("/api/config", bot.Hello)

	http.HandleFunc("/api/meetnet/safekite", bot.Safekite)
	http.HandleFunc("/api/meetnet/catalog", bot.Hello)
	http.HandleFunc("/api/meetnet/tokenexpires", bot.TokenExpires)

	//http.HandleFunc("/api/cefas/buoys", bot.CefasBuoys)
	http.HandleFunc("/api/cefas/buoys/:id", bot.Current)

	http.HandleFunc("/api/cache", bot.ListCache)
	http.HandleFunc("/api/cache/flush", bot.FlushCache)

	go func() {
		log.Println("Starting http server on port 8080")
		http.ListenAndServe(":8080", nil)
	}()
}
