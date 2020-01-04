package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func login(username string, password string) Token {
	resp, err := http.PostForm(MEETNET_API_URL+"/Token", url.Values{
		"grant_type": {"password"},
		"username":   {username},
		"password":   {password},
	})

	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if nil != err {
		log.Panic(err)
	}

	var result Token
	json.Unmarshal(body, &result)

	return result
}

func catalog(token Token) Catalog {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", MEETNET_API_URL+"/V2/catalog", nil)
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(req)

	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if nil != err {
		log.Panic(err)
	}

	var catalogResult CatalogResult
	json.Unmarshal(body, &catalogResult)

	locationMap := make(map[string]Location)
	for i := 0; i < len(catalogResult.Locations); i++ {
		locationMap[catalogResult.Locations[i].Id] = catalogResult.Locations[i]
	}

	parameterMap := make(map[string]Parameter)
	for i := 0; i < len(catalogResult.Parameters); i++ {
		parameterMap[catalogResult.Parameters[i].Id] = catalogResult.Parameters[i]
	}

	parameterTypeMap := make(map[string]ParameterType)
	for i := 0; i < len(catalogResult.ParameterTypes); i++ {
		parameterTypeMap[catalogResult.ParameterTypes[i].Id] = catalogResult.ParameterTypes[i]
	}

	availableDataMap := make(map[string]AvailableData)
	for i := 0; i < len(catalogResult.AvailableData); i++ {
		availableDataMap[catalogResult.AvailableData[i].Id] = catalogResult.AvailableData[i]
		//locationMap(catalogResult.AvailableData[i].Location)
	}

	result := Catalog{
		Locations:      locationMap,
		Parameters:     parameterMap,
		ParameterTypes: parameterTypeMap,
		AvailableData:  availableDataMap,
	}

	return result
}

func currentDataForId(token Token, ids []string) map[string]CurrentData {
	var jsonStr = ""
	if ids != nil && 0 < len(ids) {
		values := map[string][]string{"Ids": ids}
		idParam, _ :=json.Marshal(values)
		jsonStr = string(idParam)
	}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", MEETNET_API_URL+"/V2/currentData", bytes.NewBufferString(jsonStr))
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(req)

	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if nil != err {
		log.Panic(err)
	}

	var currentData []CurrentData
	json.Unmarshal(body, &currentData)

	result := make(map[string]CurrentData)
	for i := 0; i < len(currentData); i++ {
		result[currentData[i].Id] = currentData[i]
	}

	return result
}
