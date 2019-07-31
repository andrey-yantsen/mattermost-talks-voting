package web

import (
	"bytes"
	"github.com/andrey-yantsen/mattermost-talks-voting/assets/templates"
	"html/template"
	"io"
)

type Navigation struct {
	Active   bool
	Title    string
	Link     string
	Template string
}

type TemplateData struct {
	Container interface{}
	Title     string
}

type templateData struct {
	TemplateData
	Navigation []Navigation
}

var menu = []Navigation{
	{
		Active:   false,
		Title:    "Home",
		Template: "index",
		Link:     "/index",
	},
	{
		Active:   false,
		Title:    "Home2",
		Template: "index2",
		Link:     "/index2",
	},
}

func renderTemplate(templateName string, wr io.Writer, data TemplateData) error {
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
