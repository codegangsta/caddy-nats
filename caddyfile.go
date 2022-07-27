package caddynats

import (
	"strconv"

	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	httpcaddyfile.RegisterGlobalOption("nats", parseApp)
	httpcaddyfile.RegisterHandlerDirective("nats_publish", parsePublishHandler)
	httpcaddyfile.RegisterHandlerDirective("nats_request", parseRequestHandler)
	httpcaddyfile.RegisterHandlerDirective("nats_kv", parseKvHandler)
}

func parseApp(d *caddyfile.Dispenser, _ interface{}) (interface{}, error) {
	app := new(App)

	err := app.UnmarshalCaddyfile(d)

	return httpcaddyfile.App{
		Name:  "nats",
		Value: caddyconfig.JSON(app, nil),
	}, err
}

func parsePublishHandler(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var p = Publish{
		WithReply: false,
		Timeout:   publishDefaultTimeout,
	}
	err := p.UnmarshalCaddyfile(h.Dispenser)
	return p, err
}

func parseRequestHandler(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var p = Publish{
		WithReply: true,
		Timeout:   publishDefaultTimeout,
	}
	err := p.UnmarshalCaddyfile(h.Dispenser)
	return p, err
}

func parseKvHandler(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var p = KV{}
	err := p.UnmarshalCaddyfile(h.Dispenser)
	return p, err
}

func parseSubscribeHandler(d *caddyfile.Dispenser) (Subscribe, error) {
	s := Subscribe{}
	// TODO: handle errors better here
	if !d.AllArgs(&s.Subject, &s.Method, &s.URL) {
		return s, d.Err("wrong number of arguments")
	}

	return s, nil
}

func parseQueueSubscribeHandler(d *caddyfile.Dispenser) (Subscribe, error) {
	s := Subscribe{}
	// TODO: handle errors better here
	if !d.AllArgs(&s.Subject, &s.QueueGroup, &s.Method, &s.URL) {
		return s, d.Err("wrong number of arguments")
	}

	return s, nil
}

func (a *App) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			a.Context = d.Val()
		}
		if d.NextArg() {
			return d.ArgErr()
		}

		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {

			case "subscribe":
				s, err := parseSubscribeHandler(d)
				if err != nil {
					return err
				}
				jsonHandler := caddyconfig.JSONModuleObject(s, "handler", s.CaddyModule().ID.Name(), nil)
				a.HandlersRaw = append(a.HandlersRaw, jsonHandler)

			case "reply":
				s, err := parseSubscribeHandler(d)
				s.WithReply = true
				if err != nil {
					return err
				}
				jsonHandler := caddyconfig.JSONModuleObject(s, "handler", s.CaddyModule().ID.Name(), nil)
				a.HandlersRaw = append(a.HandlersRaw, jsonHandler)

			case "queue_subscribe":
				s, err := parseQueueSubscribeHandler(d)
				if err != nil {
					return err
				}
				jsonHandler := caddyconfig.JSONModuleObject(s, "handler", s.CaddyModule().ID.Name(), nil)
				a.HandlersRaw = append(a.HandlersRaw, jsonHandler)

			case "queue_reply":
				s, err := parseQueueSubscribeHandler(d)
				s.WithReply = true
				if err != nil {
					return err
				}
				jsonHandler := caddyconfig.JSONModuleObject(s, "handler", s.CaddyModule().ID.Name(), nil)
				a.HandlersRaw = append(a.HandlersRaw, jsonHandler)

			default:
				return d.Errf("unrecognized subdirective: %s", d.Val())
			}
		}
	}

	return nil
}

func (k *KV) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if !d.Args(&k.Action, &k.BucketName, &k.Key) {
			return d.Errf("Wrong argument count or unexpected line ending after '%s'", d.Val())
		}
	}

	return nil
}

func (p *Publish) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if !d.Args(&p.Subject) {
			return d.Errf("Wrong argument count or unexpected line ending after '%s'", d.Val())
		}

		for d.NextBlock(0) {
			switch d.Val() {
			case "timeout":
				if !d.NextArg() {
					return d.ArgErr()
				}
				t, err := strconv.Atoi(d.Val())
				if err != nil {
					return d.Err("timeout is not a valid integer")
				}

				p.Timeout = int64(t)
			default:
				return d.Errf("unrecognized subdirective: %s", d.Val())
			}
		}
	}

	return nil
}
