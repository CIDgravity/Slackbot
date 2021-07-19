package main

import (
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"twinQuasarAppV2/backend/customTypes"
	"twinQuasarAppV2/backend/filecoin"
)

func main() {

	// Read the configuration file
	var config customTypes.Config

	errReadConfig := cleanenv.ReadConfig("config/config.yaml", &config)

	if errReadConfig != nil {
		log.Fatal(errReadConfig)
	}

	// Initialize a session with the Filecoin Node and slack bot client
	var lotusClient, lotusCloser = filecoin.Connect(config)

	// Launch function to watchdog the fileCoin chain changes
	cacheMiners := make(map[address.Address]types.Actor)

	fmt.Println("[INFO] Launch the WatchDogs on Filecoin")
	filecoin.WaitForChange(lotusClient, cacheMiners)

	filecoin.Disconnect(lotusCloser)
}
