package http_server

import (
	"context"
	"github.com/andrey-yantsen/mattermost-talks-voting/bot"
	"log"
	"net/http"
	"strings"
	"time"
)

var Mux = http.NewServeMux()

func ListenAndServe(addr string, bot *bot.Bot) error {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "bot", bot)

		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %+v", rec)
			}
		}()

		if !strings.HasPrefix(r.URL.Path, "/cmd/") && !strings.HasPrefix(r.URL.Path, "/anonymous/") {
			authToken := r.URL.Query().Get("auth_token")
			if authToken != "" {
				http.SetCookie(w, &http.Cookie{
					Name:     "auth_token",
					Value:    authToken,
					Path:     "/",
					Expires:  time.Now().Add(time.Hour * 24),
					HttpOnly: true,
				})
				query := r.URL.Query()
				query.Del("auth_token")
				r.URL.RawQuery = query.Encode()
				http.Redirect(w, r, r.URL.String(), http.StatusPermanentRedirect)
				return
			} else {
				cookie, err := r.Cookie("auth_token")
				if err == nil && cookie.Value != "" {
					authToken = cookie.Value
					userId, channelId, exists := bot.GetDetailsFromAuthenticationToken(authToken)
					if exists {
						bot.TouchAuthenticationToken(authToken)
						ctx = context.WithValue(ctx, "user_id", userId)
						ctx = context.WithValue(ctx, "channel_id", channelId)
					} else {
						w.WriteHeader(http.StatusUnauthorized)
						w.Write([]byte("Unauthorized"))
						return
					}
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unauthorized"))
					return
				}
			}
		}

		r = r.WithContext(ctx)
		Mux.ServeHTTP(w, r)
	})
	return http.ListenAndServe(addr, handler)
}
