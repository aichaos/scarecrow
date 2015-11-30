package scarecrow

import (
	"encoding/json"
	"fmt"
	"github.com/aichaos/scarecrow/types"
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
}

func (self *Scarecrow) LoadBotsConfig() types.BotsConfig {
	config := types.BotsConfig{}

	fh, err := os.Open("config/bots.json")
	if err != nil {
		panic("Couldn't open config/bots.json; does it exist?")
	}
	defer fh.Close()

	decoder := json.NewDecoder(fh)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding bots.json:", err)
		os.Exit(1)
	}

	return config
}

// LoadUservars loads a user's RiveScript variables from disk.
func (self *Scarecrow) LoadUservars(path string) {
	if _, err := os.Stat(path); err == nil {
		fh, err := os.Open(path)
		if err != nil {
			self.Error("Unable to open profile file: %v", err)
			return
		}
		defer fh.Close()

		profile := types.UservarsConfig{}
		decoder := json.NewDecoder(fh)
		err = decoder.Decode(&profile)
		if err != nil {
			self.Error("Error decoding user profile: %s %v", path, err)
			return
		}

		self.Brain.SetUservars(profile.Username, profile.Data)
	}
}

// SaveUservars saves a user's RiveScript variables to disk.
func (self *Scarecrow) SaveUservars(username, path string) {
	vars, _ := self.Brain.GetUservars(username)

	profile := types.UservarsConfig{}
	profile.Username = username
	profile.Data = vars

	fh, err := os.Create(path)
	if err != nil {
		self.Error("Unable to create profile file: %v", err)
		return
	}
	defer fh.Close()

	encoder := json.NewEncoder(fh)
	err = encoder.Encode(profile)
	if err != nil {
		self.Error("Error encoding JSON file: %v", err)
	}
}
