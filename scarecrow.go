package scarecrow

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
	rivescript "github.com/aichaos/rivescript-go"
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
	Listeners map[string]listeners.Listener
	ListenersLock sync.RWMutex
}

func New() *Scarecrow {
	self := new(Scarecrow)
	self.Listeners = map[string]listeners.Listener{}
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

		// Make sure its ID is unique.
		if _, dupe := self.Listeners[listener.Id]; dupe {
			self.Error("Duplicate listener ID '%s'; all listeners should have a unique ID!")
			os.Exit(1)
		}

		request := make(chan types.CommunicationChannel)
		answer := make(chan types.CommunicationChannel)

		constructor, err := listeners.Create(listener.Type, listener, request, answer)
		if err != nil {
			self.Error("Unknown listener type: %s", listener.Type)
			continue
		}

		go self.ManageListener(request, answer)
		constructor.Start()
		self.Listeners[listener.Id] = constructor
	}

	self.Run()
}

// Run enters the main loop.
func (self *Scarecrow) Run() {
	for {
		time.Sleep(time.Second)
	}
}

// Shutdown shuts down all the bots.
func (self *Scarecrow) Shutdown() {
	for id, bot := range self.Listeners {
		self.Info("Send shutdown request to listener: %s", id)
		channel := bot.InputChannel()
		channel <- types.CommunicationChannel{
			Data: &types.Stop{},
		}
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

// ManageListener manages the Request/Answer channels for each listener.
func (self *Scarecrow) ManageListener(request, answer chan types.CommunicationChannel) {
	// Look for requests from the listener.
	for {
		select {
		case req := <-request:
			switch ev := req.Data.(type) {
			case *types.ReplyRequest:
				self.OnMessage(ev, answer)
			case *types.Stopped:
				self.OnStopped(ev)
			default:
				self.Error("Received an unknown event type from a listener: %v\n", ev)
			}
		}
	}
}

func (self *Scarecrow) OnMessage(req *types.ReplyRequest, res chan types.CommunicationChannel) {
	self.Log("Got reply request from %s: %s", req.Username, req.Message)
	reply := ""

	// Format the user's name to include the listener prefix, to
	// globally distinguish users on different platforms.
	uid := fmt.Sprintf("%s-%s", req.Listener, strings.ToLower(req.Username))

	// Trim their message of excess spacing.
	input := strings.Trim(req.Message, " ")

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
		} else if strings.Index(input, "!halt") == 0 {
			// !halt -- Shut the bot down.
			self.Info("Halt requested by admin user.")
			defer self.Shutdown()
			reply = "Shutting down..."
		}
	}

	if reply == "" {
		reply = self.GetReply(req.BotUsername, uid, req.Message, req.GroupChat)
	} else {
		// Log command transactions too.
		self.LogTransaction(uid, input, req.BotUsername, reply)
	}

	// Prepare an answer.
	res <- types.CommunicationChannel{
		Data: &types.ReplyAnswer{
			Username: req.Username,
			Message: reply,
		},
	}
}

// OnStopped handles when a listener informs us that they have been stopped.
func (self *Scarecrow) OnStopped(ev *types.Stopped) {
	self.ListenersLock.Lock()
	defer self.ListenersLock.Unlock()

	// Delete the listener from the stack.
	delete(self.Listeners, ev.ListenerId)

	self.Info("Listener %s has stopped. %d listeners still active.",
		ev.ListenerId,
		len(self.Listeners))

	// No remaining listeners?
	if len(self.Listeners) == 0 {
		self.Info("All listeners have stopped. Exiting the program...")
		os.Exit(0)
	}
}
