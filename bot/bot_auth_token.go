package bot

import (
	"crypto/md5"
	"fmt"
	"time"
)

type authenticationTokenData struct {
	userId    string
	channelId string
}

func (b *Bot) CreateAuthenticationToken(userId, channelId string) (token string) {
	hasher := md5.New()
	hasher.Write([]byte(userId))
	hasher.Write([]byte(channelId))
	hasher.Write([]byte(time.Now().Format(time.RFC3339Nano)))
	token = fmt.Sprintf("%x", hasher.Sum(nil))
	b.saveAuthenticationToken(token, userId, channelId)
	return
}

func (b *Bot) saveAuthenticationToken(token, userId, channelId string) {
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
	value, exists := b.cache.Get("auth_token_" + token)
	if exists {
		val := value.(authenticationTokenData)
		userId = val.userId
		channelId = val.channelId
	}
	return
}
