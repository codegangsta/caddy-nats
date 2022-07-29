package caddynats

import (
	"encoding/json"
	"fmt"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/nats-io/jsm.go/natscontext"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(App{})
}

// App connects caddy to a NATS server.
//
// NATS is a simple, secure and performant communications system for digital
// systems, services and devices.
type App struct {
	Context     string            `json:"context,omitempty"`
	HandlersRaw []json.RawMessage `json:"handle,omitempty" caddy:"namespace=nats.handlers inline_key=handler"`

	// Decoded values
	Handlers []Handler `json:"-"`

	conn   *nats.Conn
	js     nats.JetStreamContext
	logger *zap.Logger
	ctx    caddy.Context
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "nats",
		New: func() caddy.Module { return new(App) },
	}
}

// Provision sets up the app
func (app *App) Provision(ctx caddy.Context) error {
	// Set logger and Context
	app.ctx = ctx
	app.logger = ctx.Logger(app)

	// Set up handlers
	if app.HandlersRaw != nil {
		vals, err := ctx.LoadModule(app, "HandlersRaw")
		if err != nil {
			return fmt.Errorf("loading handler modules: %v", err)
		}

		for _, val := range vals.([]interface{}) {
			app.Handlers = append(app.Handlers, val.(Handler))
		}
	}

	return nil
}

func (app *App) Start() error {
	// Connect to the NATS server
	app.logger.Info("connecting via NATS context", zap.String("context", app.Context))
	conn, err := natscontext.Connect(app.Context)
	if err != nil {
		return err
	}

	js, err := conn.JetStream()
	if err != nil {
		return err
	}

	app.logger.Info("connected to NATS server",
		zap.String("url", conn.ConnectedUrlRedacted()),
		zap.String("server_id", conn.ConnectedServerId()),
		zap.String("server_name", conn.ConnectedServerName()),
		zap.String("server_version", conn.ConnectedServerVersion()),
	)
	app.conn = conn
	app.js = js

	for _, handler := range app.Handlers {
		err := handler.Subscribe(conn)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *App) Stop() error {
	defer app.conn.Close()

	app.logger.Info("closing NATS connection", zap.String("url", app.conn.ConnectedUrlRedacted()))

	for _, handler := range app.Handlers {
		err := handler.Unsubscribe(app.conn)
		if err != nil {
			return err
		}
	}

	return nil
}

// Interface guards
var (
	_ caddy.App             = (*App)(nil)
	_ caddy.Provisioner     = (*App)(nil)
	_ caddyfile.Unmarshaler = (*App)(nil)
)
