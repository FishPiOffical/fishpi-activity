package application

import (
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
)

var (
	cookieToHeader = &hook.Handler[*core.RequestEvent]{
		Id:       "cookie_to_header",
		Priority: apis.DefaultRateLimitMiddlewarePriority - 999,
		Func: func(event *core.RequestEvent) error {
			if event.Request.Header.Get("Authorization") != "" {
				return event.Next()
			}

			cookie, err := event.Request.Cookie("token")
			if err != nil {
				return event.Next()
			}
			if err = cookie.Valid(); err != nil {
				return event.Next()
			}
			event.Request.Header.Set("Authorization", cookie.Value)

			return event.Next()
		},
	}
	queryToHeader = &hook.Handler[*core.RequestEvent]{
		Id:       "query_to_header",
		Priority: apis.DefaultRateLimitMiddlewarePriority - 998,
		Func: func(event *core.RequestEvent) error {
			if event.Request.Header.Get("Authorization") != "" {
				return event.Next()
			}

			queryAuth := event.Request.URL.Query().Get("authorization")
			if queryAuth == "" {
				return event.Next()
			}
			event.Request.Header.Set("Authorization", queryAuth)

			return event.Next()
		},
	}
)
