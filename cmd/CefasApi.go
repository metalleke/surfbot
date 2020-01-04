package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Constants
const CEFAS_API_URL = "http://wavenet.cefas.co.uk/Map/GetCurrent"

//

type FeaturesCollection struct {
	Features []Feature `json:"features"`
}

type Feature struct {
	Properties Properties `json:"properties"`
}

type Properties struct {
	Id          string `json:"id"`
	WaveHeight  string `json:"WaveHeight"`
	Title       string `json:"title"`
	Temperature string `json:"Temperature"`
	Tpeak       string `json:"Tpeak"`
	Tz          string `json:"Tz"`
	Spread      string `json:"Spread"`
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
func displayCurrentWaveHeight(feature Feature) string {
	return feature.Properties.Title + ": " + feature.Properties.WaveHeight + "m"
}

func displayBuoys(features map[string]Feature) string {
	result := ""
	for _, feature := range features {
		result += displayBuoy(feature)
	}
	return result
}

func displayBuoy(feature Feature) string {
	return feature.Properties.Id + ": " + feature.Properties.Title + "\n"
}
