package web

import (
	"fmt"
	"github.com/andrey-yantsen/mattermost-talks-voting/http_server"
	"net/http"
)

func init() {
	http_server.Mux.HandleFunc("/register", HandleRegister)
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)
	channelId := r.Context().Value("channel_id").(string)

	w.Write([]byte(fmt.Sprintf("UserId: %s,\nChannelId: %s", userId, channelId)))
}
