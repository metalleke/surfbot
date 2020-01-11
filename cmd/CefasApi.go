package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Constants
const CEFAS_API_URL = "http://wavenet.cefas.co.uk/Map/GetCurrent"
const FEATURE_CACHE_KEY = "FEATURES"

//

type FeaturesCollection struct {
	Features []Feature `json:"features"`
}

type Feature struct {
	Properties Properties `json:"properties"`
}

func (t Feature) display() string {
	return t.Properties.Title + "\n" +
		"Significant wave height: " + t.Properties.WaveHeight + " m\n" +
		"Tpeak: " + t.Properties.Tpeak + " s\n" +
		"Tz: " + t.Properties.Tz + " s\n" +
		"Peak direction: " + t.Properties.Tz + " °\n" +
		"Spread: " + strconv.Itoa(t.Properties.Rotation) + " °\n" +
		"Temperature: " + t.Properties.Temperature + " °C\n"
}

func (t Feature) title() string {
	return t.Properties.Id + ": " + t.Properties.Title + "\n"
}

type Properties struct {
	Id          string `json:"id"`
	WaveHeight  string `json:"WaveHeight"`
	Title       string `json:"title"`
	Temperature string `json:"Temperature"`
	Tpeak       string `json:"Tpeak"`
	Tz          string `json:"Tz"`
	Spread      string `json:"Spread"`
	Rotation    int    `json:"rotation"`
}

// Cefas API
func getcurrent(bot *NorthSeaSurfBot) map[string]Feature {
	cachedResult, found := bot.DataCache.Remote.Get(FEATURE_CACHE_KEY)

	if found {
		return cachedResult.(map[string]Feature)
	}

	resp, err := http.Get(CEFAS_API_URL)

	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if nil != err {
		log.Panic(err)
	}

	var features FeaturesCollection
	json.Unmarshal(body, &features)

	result := make(map[string]Feature)
	for i := 0; i < len(features.Features); i++ {
		result[features.Features[i].Properties.Id] = features.Features[i]
	}

	bot.DataCache.Remote.Set(FEATURE_CACHE_KEY, result, 5*time.Minute)

	return result
}

// Displays
func displayBuoys(features map[string]Feature) string {
	result := ""
	for _, feature := range features {
		result += feature.title()
	}
	return result
}
