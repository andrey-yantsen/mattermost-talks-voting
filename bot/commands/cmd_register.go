package commands

import (
	"github.com/andrey-yantsen/mattermost-talks-voting/bot"
	"github.com/andrey-yantsen/mattermost-talks-voting/http_server"
	"github.com/andrey-yantsen/mattermost-talks-voting/storage"
	"github.com/davecgh/go-spew/spew"
	"github.com/mattermost/mattermost-server/model"
	"net/http"
)

func init() {
	http_server.Mux.HandleFunc("/cmd/register", HandleCmdRegister)
	http_server.Mux.HandleFunc("/dialog/register", HandleDialogRegister)
}

// channel_id=1uixuy6o4bnax8mnddr7u8gq1h
// &channel_name=talks-voting-bot-debug
// &command=%2Ftalks-voting-register
// &response_url=http%3A%2F%2Flocalhost%3A8065%2Fhooks%2Fcommands%2Fjopj7kpj8trbdmkra9mgpc3gxe
// &team_domain=demo
// &team_id=oybxpfkcrtgk9bz5qbbsccwmxw
// &text=
// &token=rew3b5ru6bg7bkztmndp4dq6cy
// &trigger_id=ODV1NHNlMzVjN2dhM2J0ZDNkbTh6Nm9tZHc6anJ6dTdnNGd6amRmZnBob2ZzY3dtZWJ1OGM6MTU1MTMwNjYwMjk5MDpNRVVDSVFER05Zb2hpK3dxZEE4V2l4ZGdRTzlFU0VCSEFtS3JUMnVNaHNGZjJMaHFkZ0lnWUJmZlNoRC9FVnB2MGlSdWJsUGMyaEpmK0pRLzNxYTVHbjRvUFBweld6cz0%3D
// &user_id=jrzu7g4gzjdffphofscwmebu8c
// &user_name=andrey
func HandleCmdRegister(w http.ResponseWriter, r *http.Request) {
	b := r.Context().Value("msg").(*bot.Bot)

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	form := r.Form

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := model.CommandResponse{
		TriggerId: form.Get("trigger_id"),
	}

	b.SendMsgToDebuggingChannel("```\n"+spew.Sdump(form)+"\n```", "")

	if b.IsRegistered(form.Get("channel_id")) {
		response.Text = "The channel is already registered, maybe you want to update some settings? Then call /talks-voting-update."
		response.ResponseType = "ephemeral"
		w.Write([]byte(response.ToJson()))

		return
	}

	b.SendMsgToDebuggingChannel("Got 'register' command from @"+form.Get("user_name")+" in ~"+form.Get("channel_name"), "")

	//dlg := model.Dialog{
	//	Title:       "Registering the Bot",
	//	SubmitLabel: "Register",
	//	Elements:    (&Registration{}).GetDialogElements(),
	//}
	//
	//dlgRequest := model.OpenDialogRequest{
	//	TriggerId: form.Get("trigger_id"),
	//	URL:       b.callbackUrlBase + "/dialog/register",
	//	Dialog:    dlg,
	//}
	//
	//success, dlgResponse := b.client.OpenInteractiveDialog(dlgRequest)
	//if !success {
	//	b.PrintError(dlgResponse.Error)
	//}

	response.Text = "go to somewherem"
	response.ResponseType = "ephemeral"
	// response.GotoLocation = "https://ya.ru"

	w.Write([]byte(response.ToJson()))
}

func HandleDialogRegister(w http.ResponseWriter, r *http.Request) {
	b := r.Context().Value("msg").(*bot.Bot)

	request := model.SubmitDialogRequestFromJson(r.Body)
	b.SendMsgToDebuggingChannel("```\n"+spew.Sdump(request)+"\n```", "")

	params := request.Submission
	params["channel_id"] = request.ChannelId
	params["owner_id"] = request.UserId

	reg, err := storage.LoadRegistrationFromMap(params)

	if err != nil {
		b.SendMsgToDebuggingChannel("```\n"+spew.Sdump(err)+"\n```", "")

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)

		errors := model.SubmitDialogResponse{}

		errors.Errors = make(map[string]string)

		for errorField, err := range err {
			errors.Errors[errorField] = err.Error()
		}

		w.Write([]byte(errors.ToJson()))

		return
	}

	if err := b.SaveRegistration(reg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
		return
	}
}
