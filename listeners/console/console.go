package console

import (
	"fmt"
	"github.com/jprichardson/readline-go"
	"github.com/aichaos/scarecrow/listeners"
	"github.com/aichaos/scarecrow/types"
	"os"
)

type ConsoleListener struct {
	// Channels to communicate with the parent bot.
	requestChannel chan types.CommunicationChannel
	answerChannel  chan types.CommunicationChannel

	// Configuration values for the Console listener.
	id string
	username string

	// Internal data.
	readline chan string
}

func init() {
	listeners.Register("Console", &ConsoleListener{})
}

// New creates a new Slack Listener.
func (self ConsoleListener) New(config types.ListenerConfig, request, answer chan types.CommunicationChannel) listeners.Listener {
	listener := ConsoleListener{
		id: config.Id,
		requestChannel: request,
		answerChannel: answer,
		username: config.Settings["username"],
	}

	return listener
}

func (self ConsoleListener) InputChannel() chan types.CommunicationChannel {
	return self.answerChannel
}

func (self ConsoleListener) Start() {
	self.readline = make(chan string)
	go self.ListenToConsole()
	go self.MainLoop()
}

func (self *ConsoleListener) Stop() {
	request := types.CommunicationChannel{
		Data: &types.Stopped{
			ListenerId: self.id,
		},
	}
	self.requestChannel <- request
}

func (self *ConsoleListener) ListenToConsole() {
	fmt.Printf("You> ")
	readline.ReadLine(os.Stdin, func(line string) {
		self.readline <- line
	})
}

func (self *ConsoleListener) MainLoop() {
	for {
		self.DoOneLoop()
	}
}

func (self *ConsoleListener) DoOneLoop() {
	select {
	case msg := <-self.readline:
		self.OnMessage(msg)
	case answer := <-self.answerChannel:
		switch ev := answer.Data.(type) {
		case *types.ReplyAnswer:
			self.SendMessage(ev.Username, ev.Message)
		case *types.Stop:
			self.Stop()
		}
	}
}

func (self *ConsoleListener) OnMessage(msg string) {
	request := types.CommunicationChannel{
		Data: &types.ReplyRequest{
			Listener: "CLI",
			Username: "console",
			Message: msg,
		},
	}
	self.requestChannel <- request
}

// SendMessage sends a user a response.
func (self *ConsoleListener) SendMessage(userName string, message string) {
	fmt.Printf("%s> %s\n", self.username, message)
	fmt.Printf("You> ")
}
