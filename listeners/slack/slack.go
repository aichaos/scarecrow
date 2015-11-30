package slack

import (
	"fmt"
	"strings"
	slackclient "github.com/nlopes/slack"
	"github.com/aichaos/scarecrow/types"
	"github.com/aichaos/scarecrow/listeners"
)

type SlackListener struct {
	// Channels to communicate with the parent bot.
	requestChannel chan types.ReplyRequest
	answerChannel  chan types.ReplyAnswer

	// Configuration values for the Slack listener.
	apiToken    string
	botUsername string

	// Internal data.
	api            *slackclient.Client
	rtm            *slackclient.RTM
	userChannelMap map[string]string            // Map user IDs to channel IDs
	userId2Name    map[string]string            // Map user IDs to names
	userName2Id    map[string]string            // Map user names to IDs
	userData       map[string]*slackclient.User // Full user details by ID
}

func init() {
	listeners.Register("Slack", &SlackListener{})
}

// New creates a new Slack Listener.
func (self SlackListener) New(config types.ListenerConfig, request chan types.ReplyRequest,
	response chan types.ReplyAnswer) listeners.Listener {
	listener := SlackListener{
		requestChannel: request,
		answerChannel: response,

		apiToken: config.Settings["api_token"],
		botUsername: config.Settings["username"],
	}

	listener.api = slackclient.New(listener.apiToken)

	listener.userChannelMap = map[string]string{}
	listener.userId2Name = map[string]string{}
	listener.userName2Id = map[string]string{}
	listener.userData = map[string]*slackclient.User{}

	return listener
}

func (self SlackListener) Start() {
	self.rtm = self.api.NewRTM()
	go self.rtm.ManageConnection()
	go self.MainLoop()
}

func (self *SlackListener) MainLoop() {
	for {
		self.DoOneLoop()
	}
}

func (self *SlackListener) DoOneLoop() {
	select {
	case msg := <-self.rtm.IncomingEvents:
		switch ev := msg.Data.(type) {
		case *slackclient.HelloEvent:
			// Ignore hello

		case *slackclient.ConnectedEvent:
			self.OnConnected(ev)

		case *slackclient.MessageEvent:
			self.OnMessage(ev)

		case *slackclient.RTMError:
			fmt.Printf("Slack RTM Error: %s\n", ev.Error())

		case *slackclient.InvalidAuthEvent:
			fmt.Printf("Invalid credentials\n")

		default:
			// Ignore other events.
		}
	case answer := <-self.answerChannel:
		self.SendMessage(answer.Username, answer.Message)
	}
}

// OnConnected is called when Slack connects and includes tons of data.
func (self *SlackListener) OnConnected(ev *slackclient.ConnectedEvent) {
	info := ev.Info

	// Consume the users list.
	for _, user := range info.Users {
		self.userName2Id[user.Name] = user.ID
		self.userId2Name[user.ID] = user.Name
	}
}

func (self *SlackListener) OnMessage(ev *slackclient.MessageEvent) {
	// Dissect bits of the message.
	msg := ev.Msg
	channelId := msg.Channel
	userId := msg.User
	userName := self.userId2Name[userId]
	text := msg.Text

	// Store the user ID->channel ID map.
	self.userChannelMap[userId] = channelId

	// Are we going to answer this message?
	willAnswer := false

	// Was this a public channel message or a direct message?
	if channelId[0] == 'D' {
		// Always answer DM's.
		willAnswer = true
	} else {
		// In a channel, make sure the bot's name was at-mentioned.
		atMention := fmt.Sprintf("<@%s>", self.userName2Id[self.botUsername])
		if strings.Index(text, self.botUsername) == 0 || strings.Index(text, atMention) > -1 {
			willAnswer = true
			text = strings.Replace(text, self.botUsername, "", 1)
			text = strings.Replace(text, atMention, "", -1)
			text = strings.Trim(text, ":")
		}
	}

	// Send a request for a response.
	if willAnswer {
		request := types.ReplyRequest{}
		request.BotUsername = self.botUsername
		request.Username = userName
		request.Message = text
		self.requestChannel <- request
	}
}

// SendMessage sends a user a response.
func (self *SlackListener) SendMessage(userName string, message string) {
	// Find the user ID for that name.
	userId := self.userName2Id[userName]

	// Look up the channel ID for that user ID.
	if channelId, ok := self.userChannelMap[userId]; ok {
		self.rtm.SendMessage(self.rtm.NewOutgoingMessage(message, channelId))
	}
}
