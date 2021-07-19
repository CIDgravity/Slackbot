package main

import (
	"log"
	"twinQuasarAppV2/backend/customTypes"
	"twinQuasarAppV2/backend/filecoin"
	"twinQuasarAppV2/backend/slack"
	"twinQuasarAppV2/backend/triggers"

	"github.com/ilyakaznacheev/cleanenv"
)

func main() {

	// Read the configuration file
	var config customTypes.Config

	errReadConfig := cleanenv.ReadConfig("config/config.yaml", &config)

	if errReadConfig != nil {
		log.Fatal(errReadConfig)
	}

	// Launch database triggers notification on tables that store Tipset
	triggers.InitDatabaseTriggers(config)

	// Initialize a session with the Filecoin Node and slack bot client
	var lotusClient, lotusCloser = filecoin.Connect(config)
	slack.InitSlackBot(lotusClient)

	// CLose the session with fileCoin chain
	filecoin.Disconnect(lotusCloser)
}
