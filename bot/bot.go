package bot

import (
	"fmt"
	"github.com/andrey-yantsen/mattermost-talks-voting/storage"
	"github.com/mattermost/mattermost-server/model"
	"os"
	"os/signal"
)

const (
	ApplicationName = "Mattermost Talks Voting Bot"
	ChannelLogName  = "talks-voting-bot-debug"
)

type Bot struct {
	client           *model.Client4
	botUser          *model.User
	botTeam          *model.Team
	debuggingChannel *model.Channel
	storage          *storage.Storage
	urlBase  string
	pingReceived     bool
}

func NewBot(serverUrl, accessToken, storageUrl, urlBase, teamName string, enableDebugChannel bool) *Bot {
	ret := &Bot{
		client:          model.NewAPIv4Client(serverUrl),
		urlBase: urlBase,
	}
	ret.connectStorage(storageUrl)
	ret.checkServerIsRunning()
	ret.login(accessToken)

	if enableDebugChannel {
		ret.findBotTeam(teamName)
		ret.setupGracefulShutdown()
		ret.createBotDebuggingChannelIfNeeded()
		ret.SendMsgToDebuggingChannel("_"+ApplicationName+" has **started** running_", "")
	}

	return ret
}

func (b *Bot) connectStorage(uri string) {
	s, err := storage.DbConnect(uri)

	if err != nil {
		fmt.Printf("Unable to connect to local database: %v\n", err)
		os.Exit(1)
	}

	if err := s.Migrate(); err != nil {
		fmt.Printf("Unable to apply migrations: %v\n", err)
		os.Exit(1)
	}

	b.storage = s
}

func (b *Bot) checkServerIsRunning() {
	if props, resp := b.client.GetOldClientConfig(""); resp.Error != nil {
		println("There was a problem pinging the Mattermost server.  Are you sure it's running?")
		b.PrintError(resp.Error)
		os.Exit(1)
	} else {
		println("Server detected and is running version " + props["Version"])
	}
}

func (b *Bot) login(token string) {
	b.client.MockSession(token)
	if user, resp := b.client.GetUser("me", ""); resp.Error != nil {
		println("There was a problem logging into the Mattermost server.  Are you sure ran the setup steps from the README.md?")
		b.PrintError(resp.Error)
		os.Exit(1)
	} else {
		b.botUser = user
	}
}

func (b *Bot) findBotTeam(teamName string) {
	if team, resp := b.client.GetTeamByName(teamName, ""); resp.Error != nil {
		println("We failed to get the initial load")
		println("or we do not appear to be a member of the team '" + teamName + "'")
		b.PrintError(resp.Error)
		os.Exit(1)
	} else {
		b.botTeam = team
	}
}

func (b *Bot) createBotDebuggingChannelIfNeeded() {
	if rchannel, resp := b.client.GetChannelByName(ChannelLogName, b.botTeam.Id, ""); resp.Error != nil {
		println("We failed to find the debug-log channel")
		b.PrintError(resp.Error)
	} else {
		b.debuggingChannel = rchannel
		return
	}

	// Looks like we need to create the logging channel
	channel := &model.Channel{}
	channel.Name = ChannelLogName
	channel.DisplayName = "Debugging For " + ApplicationName
	channel.Purpose = "This is used as a test channel for logging bot debug messages"
	channel.Type = model.CHANNEL_OPEN
	channel.TeamId = b.botTeam.Id
	if rchannel, resp := b.client.CreateChannel(channel); resp.Error != nil {
		println("We failed to create the channel " + ChannelLogName)
		b.PrintError(resp.Error)
	} else {
		b.debuggingChannel = rchannel
		println("Looks like this might be the first run so we've created the channel " + ChannelLogName)
	}
}

func (b *Bot) SendMsgToDebuggingChannel(msg string, replyToId string) {
	if b.debuggingChannel == nil {
		return
	}

	post := &model.Post{}
	post.ChannelId = b.debuggingChannel.Id
	post.Message = msg

	post.RootId = replyToId

	if _, resp := b.client.CreatePost(post); resp.Error != nil {
		println("We failed to send a message to the logging channel")
		b.PrintError(resp.Error)
	}
}

func (b *Bot) PrintError(err *model.AppError) {
	println("\tError Details:")
	println("\t\t" + err.Message)
	println("\t\t" + err.Id)
	println("\t\t" + err.DetailedError)
}

func (b *Bot) setupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			b.SendMsgToDebuggingChannel("_"+ApplicationName+" has **stopped** running_", "")
			os.Exit(0)
		}
	}()
}

func (b *Bot) IsRegistered(channelId string) bool {
	return b.storage.IsRegistered(channelId)
}

func (b *Bot) SaveRegistration(r *storage.Registration) error {
	return b.storage.SaveRegistration(r)
}
