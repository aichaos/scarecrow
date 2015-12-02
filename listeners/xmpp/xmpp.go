package xmpp

import (
	"crypto/tls"
	"fmt"
	"strings"
	xmppclient "github.com/aichaos/scarecrow/Godeps/_workspace/src/github.com/mattn/go-xmpp"
	"github.com/aichaos/scarecrow/listeners"
	"github.com/aichaos/scarecrow/types"
)

type XMPPListener struct {
	// Channels to communicate with the parent bot.
	requestChannel chan types.CommunicationChannel
	answerChannel  chan types.CommunicationChannel

	// Configuration values for the XMPP listener.
	id string
	username string
	port     string
	password string
	server   string

	// Internal data.
	client  *xmppclient.Client
	options xmppclient.Options
}

func init() {
	listeners.Register("XMPP", &XMPPListener{})
}

// New creates a new Slack Listener.
func (self XMPPListener) New(config types.ListenerConfig, request, answer chan types.CommunicationChannel) listeners.Listener {
	listener := XMPPListener{
		id: config.Id,
		requestChannel: request,
		answerChannel:  answer,

		server:   config.Settings["server"],
		port:     config.Settings["port"],
		username: config.Settings["username"],
		password: config.Settings["password"],
	}

	// Optional settings.
	var debug bool = config.Get("debug", "false") == "true"
	var tlsDisable bool = config.Get("notls", "false") == "true"
	var tlsNoVerify bool = config.Get("tls-no-verify", "false") == "true"
	var startTLS bool = config.Get("starttls", "false") == "true"

	// Disabling security?
	if tlsNoVerify {
		fmt.Printf("Skip TLS verify\n")
		xmppclient.DefaultConfig = tls.Config{
			ServerName: listener.server,
			InsecureSkipVerify: true,
		}
	}

	listener.options = xmppclient.Options{
		Host:          fmt.Sprintf("%s:%s", listener.server, listener.port),
		User:          listener.username,
		Password:      listener.password,
		Debug:         debug,
		Session:       true, // Use server session
		NoTLS:         tlsDisable,
		StartTLS:      startTLS,
		// Status:        "xa",
		// StatusMessage: "test status",
	}

	return listener
}

func (self XMPPListener) InputChannel() chan types.CommunicationChannel {
	return self.answerChannel
}

func (self XMPPListener) Start() {
	var err error

	self.client, err = self.options.NewClient()
	if err != nil {
		panic(fmt.Sprintf("Error connecting: %s", err))
	}

	go self.XMPPLoop()
	go self.AnswerLoop()
}

func (self *XMPPListener) Stop() {
	self.client.Close()
	self.requestChannel <- types.CommunicationChannel{
		Data: &types.Stopped{self.id},
	}
}

// XMPPLoop polls the XMPP server for incoming messages and events.
func (self *XMPPListener) XMPPLoop() {
	for {
		chat, err := self.client.Recv()
		if err != nil {
			fmt.Printf("XMPP Error: %s\n", err)
			self.Stop()
			return
		}

		switch v := chat.(type) {
		case xmppclient.Chat:
			self.OnMessage(v)
		case xmppclient.Presence:
			self.OnPresence(v)
		default:
			fmt.Printf("Unhandled XMPP event of type %s\n", v)
		}
	}
}

// AnswerLoop waits for chatbot replies to send out to the users.
func (self *XMPPListener) AnswerLoop() {
	for {
		answer := <- self.answerChannel
		switch ev := answer.Data.(type) {
		case *types.ReplyAnswer:
			self.SendMessage(ev.Username, ev.Message)
		case *types.Stop:
			self.Stop()
		}
	}
}

// OnMessage handles an incoming chat message from a user.
func (self *XMPPListener) OnMessage(v xmppclient.Chat) {
	username := v.Remote
	message := strings.Trim(v.Text, " ")

	// Remove the user's Resource from their username.
	if strings.Index(username, "/") > -1 {
		username = strings.Split(username, "/")[0]
	}

	if len(message) > 0 {
		request := types.CommunicationChannel{
			Data: &types.ReplyRequest{
				Listener: "XMPP",
				BotUsername: self.username,
				Username: username,
				Message: message,
			},
		}
		self.requestChannel <- request
	}
}

// OnPresence handles incoming presence notifications, including add requests.
func (self *XMPPListener) OnPresence(v xmppclient.Presence) {
	username := v.From

	// Remove the user's Resource from their username.
	if strings.Index(username, "/") > -1 {
		username = strings.Split(username, "/")[0]
	}

	// Handle presence types.
	if v.Type == "subscribe" {
		fmt.Printf("Subscribed by: %s\n", username)
		self.client.ApproveSubscription(username)
	}
}

// SendMessage sends a user a response.
func (self *XMPPListener) SendMessage(username string, message string) {
	self.client.Send(xmppclient.Chat{
		Remote: username,
		Type: "chat",
		Text: message,
	})
}
