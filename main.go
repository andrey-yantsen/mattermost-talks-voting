//go:generate go run assets/generate.go

package main

import (
	"flag"
	"github.com/andrey-yantsen/mattermost-talks-voting/bot"
	_ "github.com/andrey-yantsen/mattermost-talks-voting/bot/commands"
	_ "github.com/andrey-yantsen/mattermost-talks-voting/bot/web"
	"github.com/andrey-yantsen/mattermost-talks-voting/http_server"
	"log"
)

func main() {
	accessToken := flag.String("access-token", "", "Access token for the bot")
	serverUrl := flag.String("server-url", "http://localhost:8065", "Server url")
	botUrlBase := flag.String("bot-url-base", "http://localhost:8080", "Bot URL base (should be available from outside)")
	enableDebugChannel := flag.Bool("enable-debug-channel", false, "Enable debug channel")
	team := flag.String("team", "demo", "Team name in the mattermost (required when debug channel is enabled)")
	storageUri := flag.String("storage", "file:./storage/database.sqlite3?cache=shared", "SQLite database URI to use")

	flag.Parse()

	b := bot.NewBot(*serverUrl, *accessToken, *storageUri, *botUrlBase, *team, *enableDebugChannel)

	if err := http_server.ListenAndServe(":8080", b); err != nil {
		log.Fatal(err)
	}
}
