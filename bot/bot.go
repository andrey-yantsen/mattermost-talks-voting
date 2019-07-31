package bot

import (
	"context"
	"github.com/ReneKroon/ttlcache"
	"github.com/andrey-yantsen/mattermost-talks-voting/storage"
	"github.com/jinzhu/copier"
	"github.com/mattermost/mattermost-server/model"
	"log"
	"net/http"
	"net/url"
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
	urlBase          *url.URL
	pingReceived     bool
	cache            *ttlcache.Cache
}

func NewBot(serverUrl, accessToken, storageUrl, urlBase, teamName string, enableDebugChannel bool) *Bot {
	ret := &Bot{
		client: model.NewAPIv4Client(serverUrl),
		cache:  ttlcache.NewCache(),
	}
	url, err := url.Parse(urlBase)
	if err != nil {
		log.Panic(err)
	}
	ret.urlBase = url
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

func printError(err *model.AppError) {
	println("\tError Details:")
	println("\t\t" + err.Message)
	println("\t\t" + err.Id)
	println("\t\t" + err.DetailedError)
}

func (b *Bot) IsRegistered(channelId string) bool {
	return b.storage.IsRegistered(channelId)
}

func (b *Bot) SaveRegistration(r *storage.Registration) error {
	return b.storage.SaveRegistration(r)
}

func (b *Bot) CreateLink(userId, channelId, path string, params url.Values) string {
	params.Add("auth_token", b.CreateAuthenticationToken(userId, channelId))
	url := &url.URL{}
	if err := copier.Copy(url, b.urlBase); err != nil {
		log.Panic(err)
	}
	url.Path += path
	url.RawQuery = params.Encode()
	return url.String()
}

func (b *Bot) GetBotUser() *model.User {
	return b.botUser
}

func ExtractBotFromRequest(r *http.Request) *Bot {
	return ExtractBotFromContext(r.Context())
}

func ExtractBotFromContext(ctx context.Context) *Bot {
	return ctx.Value("bot").(*Bot)
}
