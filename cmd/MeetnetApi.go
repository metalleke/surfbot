package main

import (
	"fmt"
	"time"
)

// Constants

const MEETNET_API_URL = "https://api.meetnetvlaamsebanken.be"
const DEFAULT_LOCALE = "en-GB"
const SAFE_KITE_ID = "BL7WVC"

//

type Token struct {
	Expires     string `json:".expires"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	UserName    string
	Password    string
}

func (t Token) isExpired() bool {
	expires, _ := time.Parse(time.RFC1123, t.Expires)

	return expires.UTC().Add(10 * time.Second).Before(time.Now().UTC())
}

func (t Token) validate() Token {
	if t.isExpired() {
		return login(t.UserName, t.Password)
	}

	return t
}

type Location struct {
	Id            string        `json:"ID"`
	Position      string        `json:"PositionWKT"`
	Name          []Translation `json:"Name"`
	Description   []Translation `json:"Description"`
	AvailableData []AvailableData
}

func (t Location) format(locale string) string {
	return translate(t.Name, locale)
}

type Parameter struct {
	Id            string        `json:"ID"`
	Name          []Translation `json:"Name"`
	Unit          string        `json:"Unit"`
	ParameterType string        `json:"Name"`
}

type ParameterType struct {
	Id   string        `json:"ID"`
	Name []Translation `json:"Name"`
}

type AvailableData struct {
	Id              string `json:"ID"`
	Location        string `json:"Location"`
	Parameter       string `json:"Parameter"`
	CurrentInterval int    `json:"CurrentInterval"`
}

type Translation struct {
	Locale  string `json:"Culture"`
	Message string `json:"Message"`
}

type Catalog struct {
	Locations      map[string]Location
	Parameters     map[string]Parameter
	ParameterTypes map[string]ParameterType
	AvailableData  map[string]AvailableData
}

type CatalogResult struct {
	Locations      []Location      `json:"Locations"`
	Parameters     []Parameter     `json:"Parameters"`
	ParameterTypes []ParameterType `json:"ParameterTypes"`
	AvailableData  []AvailableData `json:"AvailableData"`
}

type CurrentData struct {
	Id        string  `json:"ID"`
	Timestamp string  `json:"Timestamp"`
	Value     float32 `json:"Value"`
}

func (t CurrentData) locationId() string {
	return t.Id[0:3]
}

func (t CurrentData) unitId() string {
	return t.Id[(len(t.Id) - 3):len(t.Id)]
}

type CurrentDateResult struct {
}

type SafeKite struct {
	Safe     bool
	Speed    float32
	Beaufort int
}

// Displays
func displayLocations(locations map[string]Location) string {
	result := ""
	for _, location := range locations {
		name := translate(location.Name, DEFAULT_LOCALE)
		description := translate(location.Description, DEFAULT_LOCALE)

		result += location.Id + ": " + name + " (" + description + ")\n"

	}
	return result
}

func displayCurrentData(catalog Catalog, data map[string]CurrentData) string {
	result := ""
	for _, aData := range data {
		locationId := aData.locationId()
		unitId := aData.unitId()
		location := catalog.Locations[locationId].format(DEFAULT_LOCALE)
		value := fmt.Sprintf("%.2f", aData.Value)
		unit := catalog.Parameters[unitId].Unit

		result += location + ": " + value + " " + unit + "\n"
	}
	return result
}

// Formatters
func translate(translations []Translation, locale string) string {
	for _, translation := range translations {
		if translation.Locale == locale {
			return translation.Message
		}
	}

	return ""
}

func safeToKite(beaufort int) bool {
	return beaufort < 7
}
