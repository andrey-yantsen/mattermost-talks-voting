package web

import (
	"bytes"
	"github.com/andrey-yantsen/mattermost-talks-voting/assets/templates"
	"github.com/andrey-yantsen/mattermost-talks-voting/bot"
	"github.com/mattermost/mattermost-server/model"
	"html/template"
	"io"
	"net/http"
)

type dropdownValue struct {
	Value string
	Name string
}

type dropdownValueInt struct {
	Value int
	Name string
}

type Navigation struct {
	Active   bool
	Title    string
	Link     string
	Template string
}

type TemplateData struct {
	Container interface{}
	Title     string
	Username  string
	Channel   string
}

type templateData struct {
	*TemplateData
	Navigation []Navigation
}

var menu = []Navigation{
	{
		Active:   false,
		Title:    "My channels",
		Template: "index",
		Link:     "/",
	},
}

func newTemplateDataFromRequest(r *http.Request, title string, container interface{}) *TemplateData {
	b := bot.ExtractBotFromRequest(r)
	userId := r.Context().Value("user_id").(string)
	channelId := r.Context().Value("channel_id").(string)

	userInfo := b.GetUserInfo(userId)
	channelInfo := b.GetChannelInfo(channelId)

	return newTemplateData(userInfo, channelInfo, title, container)
}

func newTemplateData(userInfo *model.User, channelInfo *model.Channel, title string, container interface{}) *TemplateData {
	username := ""
	if userInfo != nil {
		if userInfo.FirstName != "" && userInfo.LastName != "" {
			username = userInfo.LastName + " " + userInfo.FirstName
		} else {
			username = userInfo.Username
		}
	}

	channel := ""
	if channelInfo != nil {
		channel = channelInfo.DisplayName
	}

	return &TemplateData{
		Container: container,
		Title:     title,
		Username:  username,
		Channel:   channel,
	}
}

func renderTemplate(templateName string, wr io.Writer, data *TemplateData) error {
	layout, err := initTemplate(template.New("layout"), "layout")
	if err != nil {
		return err
	}
	_, err = initTemplate(layout.New("body"), templateName)
	if err != nil {
		return err
	}
	nav := make([]Navigation, len(menu))
	copy(nav, menu)
	for idx, n := range nav {
		if n.Template == templateName {
			nav[idx].Active = true
			break
		}
	}
	return layout.Execute(wr, templateData{TemplateData: data, Navigation: nav})
}

func initTemplate(template *template.Template, templateFileName string) (*template.Template, error) {
	file, err := templates.FS.Open(templateFileName + ".html")
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		return nil, err
	}
	contents := buf.String()

	_, err = template.Parse(contents)
	if err != nil {
		return nil, err
	}
	return template, nil
}
