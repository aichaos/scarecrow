package console

import (
	"fmt"
	"os"
	"github.com/jprichardson/readline-go"
	"github.com/aichaos/scarecrow/src/types"
	"github.com/aichaos/scarecrow/src/listeners"
)

type ConsoleListener struct {
	// Channels to communicate with the parent bot.
	requestChannel chan types.ReplyRequest
	answerChannel  chan types.ReplyAnswer

	// Internal data.
	username string
	readline chan string
}

func init() {
	listeners.Register("Console", &ConsoleListener{})
}

// New creates a new Slack Listener.
func (self ConsoleListener) New(config types.ListenerConfig, request chan types.ReplyRequest,
	response chan types.ReplyAnswer) listeners.Listener {
	listener := new(ConsoleListener)
	listener.requestChannel = request
	listener.answerChannel = response

	listener.username = config.Settings["username"]

	return listener
}

func (self ConsoleListener) Start() {
	self.readline = make(chan string)
	go self.ListenToConsole()
	go self.MainLoop()
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
		self.SendMessage(answer.Username, answer.Message)
	}
}

func (self *ConsoleListener) OnMessage(msg string) {
	request := types.ReplyRequest{}
	request.Username = "console"
	request.Message = msg
	self.requestChannel <- request
}

// SendMessage sends a user a response.
func (self *ConsoleListener) SendMessage(userName string, message string) {
	fmt.Printf("%s> %s\n", self.username, message)
	fmt.Printf("You> ")
}
