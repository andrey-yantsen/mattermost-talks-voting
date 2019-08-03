package web

import (
	"github.com/andrey-yantsen/mattermost-talks-voting/bot"
	"github.com/andrey-yantsen/mattermost-talks-voting/http_server"
	"log"
	"net/http"
)

func init() {
	http_server.Mux.HandleFunc("/register", HandleRegister)
}

type viewRegistrationDataContainer struct {
	ChannelName string
	TalksPerVoting int
	ViewingDays []dropdownValueInt
	SelectedViewingDay int
	ViewingTimes []dropdownValue
	SelectedViewingTime string
	Timezones []dropdownValue
	SelectedTimezone string
	MinGuestsForQuorum int
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	b := bot.ExtractBotFromRequest(r)
	userId := r.Context().Value("user_id").(string)
	channelId := r.Context().Value("channel_id").(string)

	userInfo := b.GetUserInfo(userId)
	channelInfo := b.GetChannelInfo(channelId)

	container := &viewRegistrationDataContainer{
		ChannelName: channelInfo.DisplayName,
		TalksPerVoting: 3,
		ViewingDays: getViewingDays(),
		SelectedViewingDay: 5,
		ViewingTimes: getViewingTimes(),
		SelectedViewingTime: "17:00",
		Timezones: getTimezones(),
		SelectedTimezone: "Europe/London",
		MinGuestsForQuorum: 3,
	}

	if r.Method == http.MethodPost {
		// process submitted form
	}

	data := newTemplateData(userInfo, channelInfo, "Channel registration", container)
	if err := renderTemplate("register", w, data); err != nil {
		log.Print(err)
	}
}
