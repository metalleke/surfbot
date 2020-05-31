package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Constants

const SPUIKOM_URL = "http://www.vliz.be/spuikom/widgets/tableview2.php?output=data&period=day&locations[]=Spuikomboei"

//parameter=watertemperatuurRAW

// Types

type VlizQueryResult struct {
	Dataseries []VlizDataSeries `json:"dataseries"`
}

type VlizDataSeries struct {
	Location  string      `json:"location"`
	Parameter string      `json:"parameter"`
	Units     string      `json:"units"`
	Data      [][]float32 `json:"data"`
}

type Spuikom struct {
	WaterTemperatuur string
	LuchtTemperatuur string
	Windsnelheid     string
	Windrichting     string
}

//

func (t NorthSeaSurfBot) getSpuikom() Spuikom {
	result := Spuikom{}

	luchttemperatuur := querySpuikom("luchttemperatuur")
	watertemperatuur := querySpuikom("watertemperatuurRAW")
	windrichting := querySpuikom("windrichting")
	windsnelheid := querySpuikom("windsnelheid")

	result.LuchtTemperatuur = fmt.Sprintf("%f %s", luchttemperatuur.Dataseries[0].Data[0][1], luchttemperatuur.Dataseries[0].Units)
	result.WaterTemperatuur = fmt.Sprintf("%f %s", watertemperatuur.Dataseries[0].Data[0][1], watertemperatuur.Dataseries[0].Units)
	result.Windrichting = fmt.Sprintf("%f %s", windrichting.Dataseries[0].Data[0][1], windrichting.Dataseries[0].Units)
	result.Windsnelheid = fmt.Sprintf("%f %s", windsnelheid.Dataseries[0].Data[0][1], windsnelheid.Dataseries[0].Units)

	fmt.Print(result)

	return result
}

func querySpuikom(parameter string) VlizQueryResult {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://www.vliz.be/spuikom/widgets/tableview2.php?output=data&period=day&locations[]=Spuikomboei&parameter="+parameter, nil)
	resp, err := client.Do(req)

	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if nil != err {
		log.Panic(err)
	}

	var queryResult VlizQueryResult
	json.Unmarshal(body, &queryResult)

	return queryResult
}
