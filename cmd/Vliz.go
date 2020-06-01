package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
	WaterTemperatuur float32
	LuchtTemperatuur float32
	Windsnelheid     float32
	Windrichting     float32
}

//

func (t NorthSeaSurfBot) getSpuikom() Spuikom {
	result := Spuikom{}

	luchttemperatuur := querySpuikom(t, "luchttemperatuur")
	watertemperatuur := querySpuikom(t,"watertemperatuurRAW")
	windrichting := querySpuikom(t, "windrichting")
	windsnelheid := querySpuikom(t, "windsnelheid")

	result.LuchtTemperatuur = luchttemperatuur.Dataseries[0].Data[0][1]
	result.WaterTemperatuur = watertemperatuur.Dataseries[0].Data[0][1]
	result.Windrichting = windrichting.Dataseries[0].Data[0][1]
	result.Windsnelheid = windsnelheid.Dataseries[0].Data[0][1]

	return result
}

func querySpuikom(bot NorthSeaSurfBot, parameter string) VlizQueryResult {
	cachedResult, found := bot.DataCache.Remote.Get("SPUIKOM" + parameter)

	if found {
		return cachedResult.(VlizQueryResult)
	}

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

	var result VlizQueryResult
	json.Unmarshal(body, &result)

	bot.DataCache.Remote.Set("SPUIKOM" + parameter, result, 5*time.Minute)

	return result
}
