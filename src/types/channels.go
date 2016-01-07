// Package types contains shareable types between sub-modules.
package types

////////////////////////////////////////////////////////////////////////////////
// Communication Channels for Goroutines
////////////////////////////////////////////////////////////////////////////////

type CommunicationChannel struct {
	Data interface{}
}

// ReplyRequest is a channel for a listener requesting a response for a user.
type ReplyRequest struct {
	Listener    string
	GroupChat   bool
	BotUsername string
	Username    string
	Message     string
}

// ReplyAnswer is an outgoing message to a listener to send to a user.
type ReplyAnswer struct {
	Username string
	Message  string
}

// Stop is a request for a listener to shut down.
type Stop struct{}

// Stopped is a message from a listener communicating that it has been stopped.
type Stopped struct {
	ListenerId string
}
