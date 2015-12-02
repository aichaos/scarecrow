// Package types contains shareable types between sub-modules.
package types

////////////////////////////////////////////////////////////////////////////////
// Configuration for bots.json
////////////////////////////////////////////////////////////////////////////////

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
	Id       string `json:"id"`
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
