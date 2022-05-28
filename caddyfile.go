package caddynats

import (
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	httpcaddyfile.RegisterGlobalOption("nats", parseApp)
	httpcaddyfile.RegisterHandlerDirective("nats_publish", parsePublishHandler)
	httpcaddyfile.RegisterHandlerDirective("nats_request", parseRequestHandler)

	// TODO: This is a HACK, but I couldn't find any other way to have nats_subscribe
	// and nats_reply happen in as a global directive without wrapping it in a route
	// anyway.
	httpcaddyfile.RegisterHandlerDirective("nats_subscribe", parseSubscribeHandler)
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

func parseSubscribeHandler(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var s = Subscribe{}

	err := s.UnmarshalCaddyfile(h.Dispenser)
	return s, err
}
