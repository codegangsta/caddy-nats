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
}

func parseApp(_ *caddyfile.Dispenser, _ interface{}) (interface{}, error) {
	app := new(App)

	//TODO do parsing work here

	return httpcaddyfile.App{
		Name:  "nats",
		Value: caddyconfig.JSON(app, nil),
	}, nil
}

func parsePublishHandler(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var p Publish
	err := p.UnmarshalCaddyfile(h.Dispenser)
	return p, err
}

func parseRequestHandler(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var p = Publish{WithReply: true}
	err := p.UnmarshalCaddyfile(h.Dispenser)
	return p, err
}
