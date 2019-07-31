package commands

import (
	"github.com/andrey-yantsen/mattermost-talks-voting/bot"
	"github.com/andrey-yantsen/mattermost-talks-voting/http_server"
	"github.com/mattermost/mattermost-server/model"
	"net/http"
	"net/url"
)

func init() {
	http_server.Mux.HandleFunc("/cmd/register", HandleCmdRegister)
}

func HandleCmdRegister(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b := bot.ExtractBotFromRequest(r)

	response := model.CommandResponse{
		Username:  b.GetBotUser().Username,
		TriggerId: r.Form.Get("trigger_id"),
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	form := r.Form

	if b.IsRegistered(form.Get("channel_id")) {
		response.Text = "The channel is already registered, maybe you want to update some settings? Then call /talks-voting-update."
		response.ResponseType = "ephemeral"
	} else {
		response.Text = "Please continue the registration on opened web-site."
		response.ResponseType = "ephemeral"
		response.GotoLocation = b.CreateLink(form.Get("user_id"), form.Get("channel_id"), "/register", url.Values{})
	}

	w.Write([]byte(response.ToJson()))
}
