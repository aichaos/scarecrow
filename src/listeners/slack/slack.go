package slack

import (
	"fmt"
	"regexp"
	"strings"
	slackclient "github.com/nlopes/slack"
	"github.com/aichaos/scarecrow/src/listeners"
	"github.com/aichaos/scarecrow/src/types"
)

type SlackListener struct {
	// Channels to communicate with the parent bot.
	requestChannel chan types.CommunicationChannel
	answerChannel  chan types.CommunicationChannel

	// Configuration values for the Slack listener.
	id string
	apiToken    string
	botUsername string
	team string

	// Internal data.
	api            *slackclient.Client
	rtm            *slackclient.RTM
	userChannelMap map[string]string            // Map user IDs to channel IDs
	userId2Name    map[string]string            // Map user IDs to names
	userName2Id    map[string]string            // Map user names to IDs
	userData       map[string]*slackclient.User // Full user details by ID
}

var (
	RE_MAILTO = regexp.MustCompile(`<mailto:(.+?)\|(.+?)>`)
)

func init() {
	listeners.Register("Slack", &SlackListener{})
}

// New creates a new Slack Listener.
func (self SlackListener) New(config types.ListenerConfig, request, answer chan types.CommunicationChannel) listeners.Listener {
	listener := SlackListener{
		id: config.Id,
		requestChannel: request,
		answerChannel:  answer,

		apiToken:    config.Settings["api_token"],
		botUsername: config.Settings["username"],
		team: config.Settings["team"],
	}

	listener.api = slackclient.New(listener.apiToken)

	listener.userChannelMap = map[string]string{}
	listener.userId2Name = map[string]string{}
	listener.userName2Id = map[string]string{}
	listener.userData = map[string]*slackclient.User{}

	return listener
}

func (self SlackListener) InputChannel() chan types.CommunicationChannel {
	return self.answerChannel
}

func (self SlackListener) Start() {
	self.rtm = self.api.NewRTM()
	go self.rtm.ManageConnection()
	go self.MainLoop()
}

func (self *SlackListener) Stop() {
	self.rtm.Disconnect()
	self.requestChannel <- types.CommunicationChannel{
		Data: &types.Stopped{self.id},
	}
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
		switch ev := answer.Data.(type) {
		case *types.ReplyAnswer:
			self.SendMessage(ev.Username, ev.Message)
		case *types.Stop:
			self.Stop()
		}
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

	// Ignore messages from ourself.
	if userName == self.botUsername {
		fmt.Printf("Ignore message from self: [%s] %s\n", userName, text)
		return
	}

	// Append the user's team name to the end of their nick.
	userName = fmt.Sprintf("%s@%s", userName, self.team)

	// Clean up "mailto:" links in the message.
	giveup := 0
	for strings.Index(text, "<mailto:") > -1 {
		giveup += 1
		if giveup > 50 {
			break
		}

		match := RE_MAILTO.FindStringSubmatch(text)
		if len(match) > 0 {
			pattern := fmt.Sprintf("<mailto:%s|%s>", match[1], match[2])
			text = strings.Replace(text, pattern, match[2], -1)
		}
	}


	// Store the user ID->channel ID map.
	self.userChannelMap[userId] = channelId

	willAnswer := false  // Are we going to answer this message?
	groupChat := false   // Is this a group chat message?

	// Was this a public channel message or a direct message?
	if channelId[0] == 'D' {
		// Always answer DM's.
		willAnswer = true
	} else {
		groupChat = true

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
		request := types.CommunicationChannel{
			Data: &types.ReplyRequest{
				Listener: "Slack",
				GroupChat: groupChat,
				BotUsername: self.botUsername,
				Username: userName,
				Message: text,
			},
		}
		self.requestChannel <- request
	}
}

// SendMessage sends a user a response.
func (self *SlackListener) SendMessage(userName string, message string) {
	// Strip off the team name.
	userName = strings.Split(userName, "@")[0]

	// Find the user ID for that name.
	userId := self.userName2Id[userName]

	// Look up the channel ID for that user ID.
	if channelId, ok := self.userChannelMap[userId]; ok {
		self.rtm.SendMessage(self.rtm.NewOutgoingMessage(message, channelId))
	}
}
