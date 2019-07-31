package bot

import "github.com/mattermost/mattermost-server/model"

func (b *Bot) GetUserInfo(userId string) *model.User {
	if user, response := b.client.GetUser(userId, ""); response.Error != nil {
		println("Unable to get user " + userId)
		printError(response.Error)
		return nil
	} else {
		return user
	}
}

func (b *Bot) GetChannelInfo(channelId string) *model.Channel {
	if channel, response := b.client.GetChannel(channelId, ""); response.Error != nil {
		println("Unable to get channel " + channelId)
		printError(response.Error)
		return nil
	} else {
		return channel
	}
}
