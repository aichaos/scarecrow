package scarecrow

import (
	rivescript "github.com/aichaos/rivescript-go"
	"github.com/kirsle/scarecrow/src/listeners/slack"
	"github.com/kirsle/scarecrow/src/types"
)

const (
	VERSION = "1.0.0"
)

// Type Scarecrow represents the parent object of one or more bots.
type Scarecrow struct {
	// Parameters.
	Debug bool

	// Internal structures.
	BotsConfig types.BotsConfig

	// Listeners.
	SlackListeners []*slack.SlackListener
}

type RiveScriptBot struct {
	brain *rivescript.RiveScript
}

func New() *Scarecrow {
	self := new(Scarecrow)
	self.Debug = false
	return self
}

// Start initializes and runs the bots.
func (self *Scarecrow) Start() {
	self.Info("Scarecrow version %s is starting...", VERSION)
	self.InitConfig()

	// Go through all the bots and activate them.
	for _, bot := range self.BotsConfig.Bots {
		for _, listener := range bot.Listeners {
			// Skip disabled listeners.
			if listener.Enabled == false {
				continue
			}

			// Initialize the various listener types.
			if listener.Type == "Slack" {
				request := make(chan types.ReplyRequest)
				response := make(chan types.ReplyAnswer)
				go self.ManageRequestChannel(request, response)

				obj := slack.New(listener, request, response)
				obj.Start()
				self.SlackListeners = append(self.SlackListeners, obj)
			}
		}
	}

	self.Run()
}

// Run enters the main loop.
func (self *Scarecrow) Run() {
	for {
		// Ping the Slack bots.
		for _, bot := range self.SlackListeners {
			bot.DoOneLoop()
		}
	}
}

// ManageRequestChannel manages requests for each listener.
func (self *Scarecrow) ManageRequestChannel(request chan types.ReplyRequest,
	answer chan types.ReplyAnswer) {
	// Look for requests from the listener.
	for {
		select {
		case message := <-request:
			self.Log("Got reply request from %s: %s", message.Username, message.Message)

			// Prepare an answer.
			reply := types.ReplyAnswer{}
			reply.Username = message.Username
			reply.Message = message.Message
			answer <- reply
		}
	}
}
