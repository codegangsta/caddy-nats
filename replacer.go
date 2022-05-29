package caddynats

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/nats-io/nats.go"
)

func addNATSPublishVarsToReplacer(repl *caddy.Replacer, req *http.Request, _ http.ResponseWriter, prefix string) {
	natsVars := func(key string) (any, bool) {
		if req != nil {
			switch key {
			// generated nats subject
			case "nats.subject":
				p := strings.Trim(req.URL.Path, "/")
				p = strings.TrimPrefix(p, strings.Trim(prefix, "/")+"/")
				return strings.ReplaceAll(p, "/", "."), true
			}

			// subject parts
			if strings.HasPrefix(key, natsSubjectReplPrefix) {
				idxStr := key[len(natsSubjectReplPrefix):]

				p := strings.Trim(req.URL.Path, "/")
				p = strings.TrimPrefix(p, strings.Trim(prefix, "/")+"/")
				parts := strings.Split(p, "/")

				a, b, isRange := strings.Cut(idxStr, ":")

				// range selector. n:n inclusive for both, hence the bi+1
				// TODO: This is super brute force, and I'd like to clean it up
				if isRange {
					if a == "" {
						bi, err := strconv.Atoi(b)
						if err != nil {
							return "", false
						}

						return strings.Join(parts[:bi+1], "."), true
					} else if b == "" {
						ai, err := strconv.Atoi(a)
						if err != nil {
							return "", false
						}

						return strings.Join(parts[ai:], "."), true
					} else {
						ai, err := strconv.Atoi(a)
						if err != nil {
							return "", false
						}

						bi, err := strconv.Atoi(b)
						if err != nil {
							return "", false
						}

						return strings.Join(parts[ai:bi+1], "."), true
					}
				}

				idx, err := strconv.Atoi(a)
				if err != nil {
					return "", false
				}

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

func addNatsSubscribeVarsToReplacer(repl *caddy.Replacer, msg *nats.Msg, prefix string) {
	natsVars := func(key string) (any, bool) {
		if msg != nil {
			switch key {
			// generated nats path
			case "nats.path":
				p := strings.TrimPrefix(msg.Subject, prefix+".")
				return strings.ReplaceAll(p, ".", "/"), true
			}

			// subject parts
			if strings.HasPrefix(key, natsPathReplPrefix) {
				idxStr := key[len(natsPathReplPrefix):]
				idx, err := strconv.Atoi(idxStr)
				if err != nil {
					return "", false
				}

				p := strings.TrimPrefix(msg.Subject, prefix+".")
				parts := strings.Split(p, ".")

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
var natsPathReplPrefix = "nats.path."
