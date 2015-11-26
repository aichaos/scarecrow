package scarecrow

import (
	"encoding/json"
	"fmt"
	"github.com/kirsle/scarecrow/src/types"
	"os"
)

/*
This file contains structs for config files and the necessary functionality to
load and save them.
*/

// InitConfig loads the bot's configuration files.
func (self *Scarecrow) InitConfig() {
	// Load the bots configuration.
	self.Info("Loading config: bots.json")
	self.BotsConfig = self.LoadBotsConfig()
	fmt.Printf("%v", self.BotsConfig)
}

func (self *Scarecrow) LoadBotsConfig() types.BotsConfig {
	config := types.BotsConfig{}

	file, err := os.Open("config/bots.json")
	if err != nil {
		panic("Couldn't open config/bots.json; does it exist?")
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding bots.json:", err)
		os.Exit(1)
	}

	return config
}
