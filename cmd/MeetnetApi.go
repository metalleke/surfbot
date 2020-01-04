package main

import (
	"fmt"
	"time"
)

// Constants

const MEETNET_API_URL = "https://api.meetnetvlaamsebanken.be"
const DEFAULT_LOCALE = "en-GB"

//

type Token struct {
	Expires     string `json:".expires"`
	Issued      string `json:".issued"`
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	TokenType   string `json:"token_type"`
	UserName    string
	Password    string
}

func (t Token) isExpired() bool {
	expires, _ := time.Parse(time.RFC1123, "Thu, 02 Jan 2020 20:22:36 GMT")

	return expires.Add(10 * time.Second).After(time.Now())
}

type Location struct {
	Id            string        `json:"ID"`
	Position      string        `json:"PositionWKT"`
	Name          []Translation `json:"Name"`
	Description   []Translation `json:"Description"`
	AvailableData []AvailableData
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

type CurrentDateResult struct {
}

// Helpers
func validateToken(token Token) Token {
	if token.isExpired() {
		return login(token.UserName, token.Password)
	}

	return token
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
		name := aData.Id
		locationId := name[0:3]
		unitId := name[(len(aData.Id) - 3):len(aData.Id)]
		location := translate(catalog.Locations[locationId].Name, DEFAULT_LOCALE)
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

func formatLocation(location Location, locale string) string {
	return translate(location.Name, locale)
}
