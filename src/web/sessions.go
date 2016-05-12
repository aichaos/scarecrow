package web

import (
	"crypto/rand"
	"encoding/json"
	"github.com/aichaos/scarecrow/src/log"
	"os"
)

// Session cookie management for the web server.

const KeySize = 20

// Schema for the sessions.json config file.
type sessionSchema struct {
	AuthKey    []byte `json:"authenticationSecret"`
	EncryptKey []byte `json:"encryptionSecret"`
}

var SessionConfig sessionSchema

// On first initialization, generate suitably random keys for session security.
func init() {
	if _, err := os.Stat("config/sessions.json"); os.IsNotExist(err) {
		// No session found yet, generate new keys.
		initSessionKeys()
	} else {
		// Load the existing keys.
		loadSessionKeys()
	}
}

// initSessionKeys creates new session security keys and saves them to disk.
func initSessionKeys() {
	log.Debug("No session keys found, generating new keys")

	// Generate new keys from /dev/urandom.
	authKey := make([]byte, KeySize)
	encryptKey := make([]byte, KeySize)
	_, _ = rand.Read(authKey)
	_, _ = rand.Read(encryptKey)

	// Put them in the session config.
	SessionConfig = sessionSchema{
		AuthKey:    authKey,
		EncryptKey: encryptKey,
	}

	fh, err := os.Create("config/sessions.json")
	if err != nil {
		log.Error("Unable to create sessions config: %v", err)
		os.Exit(1)
	}
	defer fh.Close()

	encoder := json.NewEncoder(fh)
	err = encoder.Encode(&SessionConfig)
	if err != nil {
		log.Error("Error encoding sessions.json: %v", err)
		os.Exit(1)
	}
}

func loadSessionKeys() {
	log.Debug("Loading session keys from config/sessions.json")
	SessionConfig = sessionSchema{}

	fh, err := os.Open("config/sessions.json")
	if err != nil {
		log.Error("Unable to load sessions config: %v", err)
		os.Exit(1)
	}
	defer fh.Close()

	decoder := json.NewDecoder(fh)
	err = decoder.Decode(&SessionConfig)
	if err != nil {
		log.Error("Error decoding sessions.json: %v", err)
		os.Exit(1)
	}
}
