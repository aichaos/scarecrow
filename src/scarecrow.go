package scarecrow

import (
	rivescript "github.com/aichaos/rivescript-go"
	"github.com/aichaos/scarecrow/src/listeners"
	"github.com/aichaos/scarecrow/src/types"
	"strings"
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
	Listeners []listeners.Listener
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
		self.Info("Setting up %s listener...", listener.Type)
		request := make(chan types.ReplyRequest)
		response := make(chan types.ReplyAnswer)
		go self.ManageRequestChannel(request, response)

		constructor, err := listeners.Create(listener.Type, listener, request, response)
		if err != nil {
			self.Error("Unknown listener type: %s", listener.Type)
			continue
		}

		constructor.Start()
		self.Listeners = append(self.Listeners, constructor)
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

			input := strings.Trim(message.Message, " ")
			var reply string

			// Handle commands (TODO: admin rights and such)
			if strings.Index(input, "!reload") == 0 {
				self.InitBrain()
				reply = "Brain reloaded!"
			} else {
				reply = self.GetReply(message.BotUsername, message.Username, message.Message)
			}

			// Prepare an answer.
			outgoing := types.ReplyAnswer{}
			outgoing.Username = message.Username
			outgoing.Message = reply
			answer <- outgoing
		}
	}
}
