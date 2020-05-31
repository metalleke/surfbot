package main

import "testing"

func TestGetSpuikom(t *testing.T) {
	bot := NorthSeaSurfBot{
		Config:    Config{},
		DataCache: DataCache{},
	}

	bot.getSpuikom()

}
