package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Constants
const CEFAS_API_URL = "http://wavenet.cefas.co.uk/Map/GetCurrent"

// Types
type FeaturesCollection struct {
	Features []Feature `json:"features"`
}

type Feature struct {
	Properties Properties `json:"properties"`
}

func (t Feature) display() string {
	return t.Properties.Title + "\n" +
		"Significant wave height: " + t.Properties.WaveHeight + " m\n" +
		"Tpeak: " + t.Properties.WaveHeight + " s\n" +
		"Tz: " + t.Properties.WaveHeight + " s\n" +
		"Peak direction: " + t.Properties.WaveHeight + " °\n" +
		"Spread: " + t.Properties.WaveHeight + " °\n" +
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

// Retrieve data
func getcurrent() map[string]Feature {
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
