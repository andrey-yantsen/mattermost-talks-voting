package commands

import (
	"github.com/andrey-yantsen/mattermost-talks-voting/http_server"
	"net/http"
)

func init() {
	http_server.Mux.HandleFunc("/cmd/add", HandleCmdAdd)
}

func HandleCmdAdd(w http.ResponseWriter, r *http.Request) {

}
