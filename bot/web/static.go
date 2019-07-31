package web

import (
	"github.com/andrey-yantsen/mattermost-talks-voting/assets/static"
	"github.com/andrey-yantsen/mattermost-talks-voting/http_server"
	"net/http"
)

func init() {
	http.Handle()
	http_server.Mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(static.FS)))
}
