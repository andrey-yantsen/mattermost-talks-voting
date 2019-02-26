// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package main

import (
	"flag"
	"github.com/mattermost/mattermost-server/model"
	"os"
	"os/signal"
	"regexp"
	"strings"
)

const (
	APPLICATION_NAME = "Mattermost Talks Voting Bot"
	CHANNEL_LOG_NAME = "talks-voting-bot-debug"
)

var client *model.Client4
var webSocketClient *model.WebSocketClient

var botUser *model.User
var botTeam *model.Team
var debuggingChannel *model.Channel

func main() {
	userEmail := flag.String("user-email", "talks-voting@example.com", "Email set for the bot")
	password := flag.String("password", "password", "Password for the bot")
	serverUrl := flag.String("server-url", "http://localhost:8065", "Server url")
	team := flag.String("team", "demo", "Team name in the mattermost")
	createDebugChannel := flag.Bool("create-debug-channel", false, "Whenever to create a debug channel")

	flag.Parse()

	SetupGracefulShutdown()

	client = model.NewAPIv4Client(*serverUrl)

	MakeSureServerIsRunning()
	LoginAsTheBotUser(*userEmail, *password)
	FindBotTeam(*team)

	if *createDebugChannel {
		CreateBotDebuggingChannelIfNeeded()
		SendMsgToDebuggingChannel("_"+APPLICATION_NAME+" has **started** running_", "")
	}

	wsUrl := strings.Replace(strings.Replace(*serverUrl, "http://", "ws://", 1), "https://", "ws://", 1)
	webSocketClient, err := model.NewWebSocketClient4(wsUrl, client.AuthToken)
	if err != nil {
		println("We failed to connect to the web socket")
		PrintError(err)
	}

	webSocketClient.Listen()

	go func() {
		for {
			select {
			case resp := <-webSocketClient.EventChannel:
				HandleWebSocketResponse(resp)
			}
		}
	}()

	// You can block forever with
	select {}
}

func MakeSureServerIsRunning() {
	if props, resp := client.GetOldClientConfig(""); resp.Error != nil {
		println("There was a problem pinging the Mattermost server.  Are you sure it's running?")
		PrintError(resp.Error)
		os.Exit(1)
	} else {
		println("Server detected and is running version " + props["Version"])
	}
}

func LoginAsTheBotUser(email, password string) {
	if user, resp := client.Login(email, password); resp.Error != nil {
		println("There was a problem logging into the Mattermost server.  Are you sure ran the setup steps from the README.md?")
		PrintError(resp.Error)
		os.Exit(1)
	} else {
		botUser = user
	}
}

func FindBotTeam(teamName string) {
	if team, resp := client.GetTeamByName(teamName, ""); resp.Error != nil {
		println("We failed to get the initial load")
		println("or we do not appear to be a member of the team '" + teamName + "'")
		PrintError(resp.Error)
		os.Exit(1)
	} else {
		botTeam = team
	}
}

func CreateBotDebuggingChannelIfNeeded() {
	if rchannel, resp := client.GetChannelByName(CHANNEL_LOG_NAME, botTeam.Id, ""); resp.Error != nil {
		println("We failed to get the channels")
		PrintError(resp.Error)
	} else {
		debuggingChannel = rchannel
		return
	}

	// Looks like we need to create the logging channel
	channel := &model.Channel{}
	channel.Name = CHANNEL_LOG_NAME
	channel.DisplayName = "Debugging For " + APPLICATION_NAME
	channel.Purpose = "This is used as a test channel for logging bot debug messages"
	channel.Type = model.CHANNEL_OPEN
	channel.TeamId = botTeam.Id
	if rchannel, resp := client.CreateChannel(channel); resp.Error != nil {
		println("We failed to create the channel " + CHANNEL_LOG_NAME)
		PrintError(resp.Error)
	} else {
		debuggingChannel = rchannel
		println("Looks like this might be the first run so we've created the channel " + CHANNEL_LOG_NAME)
	}
}

func SendMsgToDebuggingChannel(msg string, replyToId string) {
	if debuggingChannel == nil {
		return
	}

	post := &model.Post{}
	post.ChannelId = debuggingChannel.Id
	post.Message = msg

	post.RootId = replyToId

	if _, resp := client.CreatePost(post); resp.Error != nil {
		println("We failed to send a message to the logging channel")
		PrintError(resp.Error)
	}
}

func HandleWebSocketResponse(event *model.WebSocketEvent) {
	if debuggingChannel != nil && event.Broadcast.ChannelId == debuggingChannel.Id {
		HandleMsgFromDebuggingChannel(event)
	}
}

func HandleMsgFromDebuggingChannel(event *model.WebSocketEvent) {
	if event.Event != model.WEBSOCKET_EVENT_POSTED {
		return
	}

	println("responding to debugging channel msg")

	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
	if post != nil {
		if post.UserId == botUser.Id {
			return
		}

		if matched, _ := regexp.MatchString(`(?:^|\W)alive(?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("Yes I'm running", post.Id)
			return
		}

		if matched, _ := regexp.MatchString(`(?:^|\W)up(?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("Yes I'm running", post.Id)
			return
		}

		if matched, _ := regexp.MatchString(`(?:^|\W)running(?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("Yes I'm running", post.Id)
			return
		}

		if matched, _ := regexp.MatchString(`(?:^|\W)hello(?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("Yes I'm running", post.Id)
			return
		}
	}

	SendMsgToDebuggingChannel("I did not understand you!", post.Id)
}

func PrintError(err *model.AppError) {
	println("\tError Details:")
	println("\t\t" + err.Message)
	println("\t\t" + err.Id)
	println("\t\t" + err.DetailedError)
}

func SetupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			if webSocketClient != nil {
				webSocketClient.Close()
			}

			SendMsgToDebuggingChannel("_"+APPLICATION_NAME+" has **stopped** running_", "")
			os.Exit(0)
		}
	}()
}
