// Package types contains shareable types between sub-modules.
package types

////////////////////////////////////////////////////////////////////////////////
// Communication Channels for Goroutines
////////////////////////////////////////////////////////////////////////////////

// ReplyRequest is a channel for a listener requesting a response for a user.
type ReplyRequest struct {
	Username string
	Message  string
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
	Bots []BotConfig
}

type BotConfig struct {
	Name      string           `json:"name"`
	Brain     BrainConfig      `json:"brain"`
	Listeners []ListenerConfig `json:"listeners"`
}

type BrainConfig struct {
	Backend string `json:"backend"`
	Replies string `json:"replies"`
}

type ListenerConfig struct {
	Type     string `json:"type"`
	APIToken string `json:"api_token"`
	Enabled  bool   `json:"enabled"`
	Username string `json:"username"`
}
