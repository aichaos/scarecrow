package scarecrow

import (
	rivescript "github.com/aichaos/rivescript-go"
	"github.com/aichaos/scarecrow/src/listeners/console"
	"github.com/aichaos/scarecrow/src/listeners/slack"
	"github.com/aichaos/scarecrow/src/types"
	"time"
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
	Brain      *rivescript.RiveScript

	// Listeners.
	SlackListeners   []*slack.SlackListener
	ConsoleListeners []*console.ConsoleListener
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
	self.InitBrain()
	MakeDirectory("./users")

	// Go through all the bots and activate them.
	for _, listener := range self.BotsConfig.Listeners {
		// Skip disabled listeners.
		if listener.Enabled == false {
			continue
		}

		// Initialize the various listener types.
		if listener.Type == "Console" {
			self.Info("Setting up Console listener...")
			request := make(chan types.ReplyRequest)
			response := make(chan types.ReplyAnswer)
			go self.ManageRequestChannel(request, response)

			obj := console.New(listener, request, response)
			obj.Start()
			self.ConsoleListeners = append(self.ConsoleListeners, obj)
		} else if listener.Type == "Slack" {
			self.Info("Setting up Slack listener...")
			request := make(chan types.ReplyRequest)
			response := make(chan types.ReplyAnswer)
			go self.ManageRequestChannel(request, response)

			obj := slack.New(listener, request, response)
			obj.Start()
			self.SlackListeners = append(self.SlackListeners, obj)
		} else {
			self.Warn("Unknown listener type %s", listener.Type)
		}
	}

	self.Run()
}

// Run enters the main loop.
func (self *Scarecrow) Run() {
	for {
		time.Sleep(time.Second)
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
			reply := self.GetReply(message.BotUsername, message.Username, message.Message)

			// Prepare an answer.
			outgoing := types.ReplyAnswer{}
			outgoing.Username = message.Username
			outgoing.Message = reply
			answer <- outgoing
		}
	}
}
