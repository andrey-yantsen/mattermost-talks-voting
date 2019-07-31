package web

import (
	"github.com/andrey-yantsen/mattermost-talks-voting/bot"
	"github.com/andrey-yantsen/mattermost-talks-voting/http_server"
	"net/http"
)

func init() {
	http_server.Mux.HandleFunc("/register", HandleRegister)
}

type viewRegistrationDataContainer struct {
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	b := bot.ExtractBotFromRequest(r)
	userId := r.Context().Value("user_id").(string)
	channelId := r.Context().Value("channel_id").(string)

	userInfo := b.GetUserInfo(userId)
	channelInfo := b.GetChannelInfo(channelId)

	container := &viewRegistrationDataContainer{}

	if r.Method == http.MethodPost {
		// process submitted form
	}

	data := newTemplateData(userInfo, channelInfo, "Channel registration", container)
	renderTemplate("register", w, data)
}
