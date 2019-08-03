package bot

import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"
)

type authenticationTokenData struct {
	userId    string
	channelId string
}

func (b *Bot) CreateAuthenticationToken(userId, channelId string) (token string) {
	token = fmt.Sprintf("%s:%s", userId, channelId)
	return
	hasher := md5.New()
	hasher.Write([]byte(userId))
	hasher.Write([]byte(channelId))
	hasher.Write([]byte(time.Now().Format(time.RFC3339Nano)))
	token = fmt.Sprintf("%x", hasher.Sum(nil))
	b.saveAuthenticationToken(token, userId, channelId)
	return
}

func (b *Bot) saveAuthenticationToken(token, userId, channelId string) {
	return
	b.cache.SetWithTTL("auth_token_"+token, authenticationTokenData{userId, channelId}, 15*time.Minute)
}

func (b *Bot) TouchAuthenticationToken(token string) bool {
	if userId, channelId, exists := b.GetDetailsFromAuthenticationToken(token); exists {
		b.saveAuthenticationToken(token, userId, channelId)
		return true
	} else {
		return false
	}
}

func (b *Bot) GetDetailsFromAuthenticationToken(token string) (userId, channelId string, exists bool) {
	data := strings.Split(token, ":")
	userId = data[0]
	channelId = data[1]
	exists = true
	return
	value, exists := b.cache.Get("auth_token_" + token)
	if exists {
		val := value.(authenticationTokenData)
		userId = val.userId
		channelId = val.channelId
	}
	return
}
