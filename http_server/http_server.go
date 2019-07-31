package http_server

import (
	"context"
	"github.com/andrey-yantsen/mattermost-talks-voting/bot"
	"log"
	"net/http"
)

var Mux = http.NewServeMux()

func ListenAndServe(addr string, bot *bot.Bot) error {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(context.Background(), "bot", bot)
		r = r.WithContext(ctx)

		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %+v", rec)
			}
		}()

		Mux.ServeHTTP(w, r)
	})
	return http.ListenAndServe(addr, handler)
}
