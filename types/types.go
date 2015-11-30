// Package types contains shareable types between sub-modules.
package types

////////////////////////////////////////////////////////////////////////////////
// Communication Channels for Goroutines
////////////////////////////////////////////////////////////////////////////////

// ReplyRequest is a channel for a listener requesting a response for a user.
type ReplyRequest struct {
	BotUsername string
	Username    string
	Message     string
}

// ReplyAnswer is an outgoing message to a listener to send to a user.
type ReplyAnswer struct {
	Username string
	Message  string
}

////////////////////////////////////////////////////////////////////////////////
// Configuration File Types
////////////////////////////////////////////////////////////////////////////////

/****************************/
/***** config/bots.json *****/
/****************************/

type BotsConfig struct {
	Personality PersonalityConfig `json:"personality"`
	Listeners   []ListenerConfig  `json:"listeners"`
}

type PersonalityConfig struct {
	Name  string      `json:"name"`
	Brain BrainConfig `json:"brain"`
}

type BrainConfig struct {
	Backend string `json:"backend"`
	Replies string `json:"replies"`
}

type ListenerConfig struct {
	Type     string `json:"type"`
	Enabled  bool   `json:"enabled"`
	Settings map[string]string `json:"settings"`
}

// Get safely gets an optional config key or falls back to a default value.
func (self *ListenerConfig) Get(name, fallback string) string {
	if value, ok := self.Settings[name]; ok {
		return value
	}
	return fallback
}

/********************************/
/***** RiveScript User Data *****/
/********************************/

type UservarsConfig struct {
	Username string            `json:"username"`
	Data     map[string]string `json:"data"`
}
