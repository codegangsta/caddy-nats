package caddynats

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/caddyserver/caddy/v2"
)

func addNATSVarsToReplacer(repl *caddy.Replacer, req *http.Request, _ http.ResponseWriter) {
	natsVars := func(key string) (any, bool) {
		if req != nil {
			switch key {
			// generated nats subject
			case "nats.subject":
				p := strings.Trim(req.URL.Path, "/")
				return strings.ReplaceAll(p, "/", "."), true
			}

			// subject parts
			if strings.HasPrefix(key, natsSubjectReplPrefix) {
				idxStr := key[len(natsSubjectReplPrefix):]
				idx, err := strconv.Atoi(idxStr)
				if err != nil {
					return "", false
				}

				p := strings.Trim(req.URL.Path, "/")
				parts := strings.Split(p, "/")

				if idx < 0 {
					return "", false
				}

				if idx >= len(parts) {
					return "", true
				}

				return parts[idx], true
			}

		}

		return nil, false
	}

	repl.Map(natsVars)
}

var natsSubjectReplPrefix = "nats.subject."
