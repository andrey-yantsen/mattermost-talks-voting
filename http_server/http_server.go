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

func outerHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %+v", rec)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func injectBot(bot *bot.Bot, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "bot", bot)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func replaceAuthTokenParamWithCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		}
		next.ServeHTTP(w, r)
	})
}

func parseAuthCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		ctx := r.Context()
		if err == nil && cookie.Value != "" {
			bot := bot.ExtractBotFromContext(r)
			authToken := cookie.Value
			userId, channelId, exists := bot.GetDetailsFromAuthenticationToken(authToken)
			if exists {
				bot.TouchAuthenticationToken(authToken)
				ctx = context.WithValue(r.Context(), "user_id", userId)
				ctx = context.WithValue(ctx, "channel_id", channelId)
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func denyUnauthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/cmd/") && !strings.HasPrefix(r.URL.Path, "/anonymous/") {
			ctx := r.Context()
			userId := ctx.Value("user_id")
			channelId := ctx.Value("channel_id")
			if userId == nil || channelId == nil || userId.(string) == "" || channelId.(string) == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func ListenAndServe(addr string, bot *bot.Bot) error {
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Mux.ServeHTTP(w, r)
	})
	return http.ListenAndServe(addr, outerHandler(injectBot(bot, replaceAuthTokenParamWithCookie(parseAuthCookie(denyUnauthorized(finalHandler))))))
}
