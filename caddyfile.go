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
				s := Subscribe{}
				d.Prev()
				err := s.UnmarshalCaddyfile(d)
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

func (p *Publish) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if !d.Args(&p.Subject) {
			return d.Errf("Wrong argument count or unexpected line ending after '%s'", d.Val())
		}

		for d.NextBlock(0) {
			switch d.Val() {
			case "prefix":
				if p.Prefix != "" {
					return d.Err("prefix already specified")
				}
				if !d.NextArg() {
					return d.ArgErr()
				}

				p.Prefix = d.Val()
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

func (s *Subscribe) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		// TODO better error handling
		d.Args(&s.Subject, &s.Method, &s.Path)
	}

	return nil
}
