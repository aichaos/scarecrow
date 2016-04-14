// Package db provides the database driver for Scarecrow.
package db

import (
	"encoding/json"
	"os"
	"github.com/aichaos/scarecrow/src/log"
	"github.com/aichaos/scarecrow/src/models"
	"github.com/aichaos/scarecrow/src/types"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Driver gorm.DB
	Config types.DBConfig
	Ready  bool
}

var instance *DB

// New creates a new database handler.
func New(conf types.DBConfig) *DB {
	if instance != nil {
		return instance
	}

	db := new(DB)
	db.Config = conf
	db.Ready  = false

	// If we have valid configuration, attempt to connect.
	if len(conf.Type) > 0 && len(conf.ConnString) > 0 {
		log.Debug("Attempting to connect to the database...")
		db.Connect()
	}

	instance = db
	return instance
}

// GetInstance returns a previously created DB instance.
func GetInstance() *DB {
	return instance
}

// Setup attempts to connect to the database. If no database configuration
// exists, this will return false.
func (self *DB) Connect() (bool, error) {
	db, err := gorm.Open(self.Config.Type, self.Config.ConnString)
	if err != nil {
		log.Error("Failed to connect to database: %s", err)
		return false, err
	}

	self.Driver = db
	self.Ready  = true

	self.CreateTables()

	return true, nil
}

// CreateTables makes sure all the database tables exist.
func (self *DB) CreateTables() {
	self.Driver.CreateTable(&models.User{})
	self.Driver.CreateTable(&models.Bot{})
	self.Driver.CreateTable(&models.Listener{})
	self.Driver.CreateTable(&models.ListenerSetting{})
}

// LoadConfig loads the config/db.json config file.
func LoadConfig() types.DBConfig {
	config := types.DBConfig{}

	fh, err := os.Open("config/db.json")
	if err != nil {
		log.Warn("No database configuration found. Please go to the web app to set up your bot.")
		return config
	}
	defer fh.Close()

	decoder := json.NewDecoder(fh)
	err = decoder.Decode(&config)
	if err != nil {
		log.Error("Error decoding db.json:", err)
		os.Exit(1)
	}

	return config
}

// SaveConfig saves the database configuration to disk.
func (self *DB) SaveConfig() {
	cfg := self.Config

	fh, err := os.Create("config/db.json")
	if err != nil {
		log.Error("Unable to create database config file: %v", err)
		return
	}
	defer fh.Close()

	encoder := json.NewEncoder(fh)
	err = encoder.Encode(cfg)
	if err != nil {
		log.Error("Error encoding JSON file: %v", err)
	}
}
