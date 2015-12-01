// Package types contains shareable types between sub-modules.
package types

////////////////////////////////////////////////////////////////////////////////
// Communication Channels for Goroutines
////////////////////////////////////////////////////////////////////////////////

// ReplyRequest is a channel for a listener requesting a response for a user.
type ReplyRequest struct {
	Listener string
	BotUsername string
	Username    string
	Message     string
}

// ReplyAnswer is an outgoing message to a listener to send to a user.
type ReplyAnswer struct {
	Username string
	Message  string
}
