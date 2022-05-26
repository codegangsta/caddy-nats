package caddynats

import (
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
)

func addNATSVarsToReplacer(repl *caddy.Replacer, req *http.Request, _ http.ResponseWriter) {
	natsVars := func(key string) (any, bool) {
		if req != nil {
			switch key {
			case "nats.subject":
				p := strings.Trim(req.URL.Path, "/")
				return strings.ReplaceAll(p, "/", "."), true
			}
		}

		return nil, false
	}

	repl.Map(natsVars)
}
