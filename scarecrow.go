package scarecrow

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	rivescript "github.com/aichaos/scarecrow/Godeps/_workspace/src/github.com/aichaos/rivescript-go"
	"github.com/aichaos/scarecrow/listeners"
	"github.com/aichaos/scarecrow/types"
)

const (
	VERSION = "1.0.0"
)

var (
	RE_OP = regexp.MustCompile(`^!op ([A-Za-z0-9\.@\-_]+?)$`)
	RE_DEOP = regexp.MustCompile(`^!deop ([A-Za-z0-9\.@\-_]+?)$`)
)

// Type Scarecrow represents the parent object of one or more bots.
type Scarecrow struct {
	// Parameters.
	Debug bool

	// Internal structures.
	AdminsConfig types.AdminsConfig
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

// IsAdmin returns whether a user ID is an admin user or not.
func (self *Scarecrow) IsAdmin(username string) bool {
	for _, user := range self.AdminsConfig.Admins {
		if user == username {
			return true
		}
	}
	return false
}

// ManageRequestChannel manages requests for each listener.
func (self *Scarecrow) ManageRequestChannel(request chan types.ReplyRequest,
	answer chan types.ReplyAnswer) {
	// Look for requests from the listener.
	for {
		select {
		case message := <-request:
			self.Log("Got reply request from %s: %s", message.Username, message.Message)
			reply := ""

			// Format the user's name to include the listener prefix, to
			// globally distinguish users on different platforms.
			uid := fmt.Sprintf("%s-%s", message.Listener, message.Username)

			// Trim their message of excess spacing.
			input := strings.Trim(message.Message, " ")

			// Handle commands (TODO: admin rights and such)
			if self.IsAdmin(uid) {
				if strings.Index(input, "!reload") == 0 {
					// !reload -- Reload the RiveScript brain.
					self.InitBrain()
					reply = "Brain reloaded!"
				} else if strings.Index(input, "!op") == 0 {
					// !op -- Add a user as an admin.
					match := RE_OP.FindStringSubmatch(input)
					if len(match) > 0 {
						opName := match[1]
						self.AdminsConfig.Admins = append(self.AdminsConfig.Admins, opName)
						self.SaveAdminsConfig(self.AdminsConfig)
						reply = fmt.Sprintf("%s added to the admins list.", opName)
					} else {
						self.Warn("Syntax error parsing command: %s", input)
						reply = "Syntax error."
					}
				} else if strings.Index(input, "!deop") == 0 {
					// !deop -- Remove a user as an admin.
					match := RE_DEOP.FindStringSubmatch(input)
					if len(match) > 0 {
						opName := match[1]

						// Remove them from the list.
						newAdmins := []string{}
						for _, name := range self.AdminsConfig.Admins {
							if name != opName {
								newAdmins = append(newAdmins, name)
							}
						}
						self.AdminsConfig.Admins = newAdmins
						self.SaveAdminsConfig(self.AdminsConfig)

						reply = fmt.Sprintf("%s removed from the admins list.", opName)
					} else {
						self.Warn("Syntax error parsing command: %s", input)
						reply = "Syntax error."
					}
				}
			}

			if reply == "" {
				reply = self.GetReply(message.BotUsername, uid, message.Message)
			} else {
				// Log command transactions too.
				self.LogTransaction(uid, input, message.BotUsername, reply)
			}

			// Prepare an answer.
			answer <- types.ReplyAnswer{
				Username: message.Username,
				Message: reply,
			}
		}
	}
}
