package caddynats

import (
	"github.com/caddyserver/caddy/v2"
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
	conn *nats.Conn

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
	return nil
}

func (app *App) Start() error {
	// Connect to the NATS server
	//TODO Configure URL and options
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		app.logger.Error(err.Error())
	}

	app.logger.Info("Connected to NATS server", zap.String("url", conn.ConnectedUrlRedacted()))
	app.conn = conn

	return nil
}

func (app *App) Stop() error {
	app.logger.Info("Closing NATS connection", zap.String("url", app.conn.ConnectedUrlRedacted()))
	// TODO: Do we need to drain here?
	app.conn.Close()
	return nil
}

// Interface guards
var (
	_ caddy.App         = (*App)(nil)
	_ caddy.Provisioner = (*App)(nil)
)
