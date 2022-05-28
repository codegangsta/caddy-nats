package caddynats

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Subscribe struct {
	Subject string `json:"subject,omitempty"`
	Method  string `json:"method,omitempty"`
	Path    string `json:"path,omitempty"`

	sub     *nats.Subscription
	logger  *zap.Logger
	app     *App
	httpApp *caddyhttp.App
}

func init() {
	caddy.RegisterModule(Subscribe{})
}

func (Subscribe) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.nats_subscribe",
		New: func() caddy.Module { return new(Subscribe) },
	}
}

func (s *Subscribe) Provision(ctx caddy.Context) error {
	s.logger = ctx.Logger(s)

	natsAppIface, err := ctx.App("nats")
	if err != nil {
		return fmt.Errorf("getting NATS app: %v. Make sure NATS is configured in global options", err)
	}
	s.app = natsAppIface.(*App)

	s.logger.Info("Subscribing to NATS subject", zap.String("subject", s.Subject))
	sub, err := s.app.conn.Subscribe(s.Subject, s.handler)
	s.sub = sub

	return err
}

func (s Subscribe) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Do nothing, since this is a subscriber
	return next.ServeHTTP(w, r)
}

func (s *Subscribe) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		// TODO better error handling
		d.Args(&s.Subject, &s.Method, &s.Path)
	}

	return nil
}

func (s *Subscribe) handler(msg *nats.Msg) {
	s.logger.Info("Handling message on subject", zap.String("subject", msg.Subject))

	// Look up the http app
	httpAppIface, err := s.app.ctx.App("http")
	if err != nil {
		s.logger.Error("http app not loaded", zap.String("subject", msg.Subject))
		return
	}
	httpApp := httpAppIface.(*caddyhttp.App)

	// TODO: Support multiple servers, for now just pick the first one
	server := httpApp.Servers["srv0"]

	req, err := http.NewRequest(s.Method, s.Path, bytes.NewBuffer(msg.Data))
	if err != nil {
		// TODO: don't panic
		panic(err)
	}

	// TODO: Only use this for replies, otherwise use a no-op recorder
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	msg.Respond(rec.Body.Bytes())
}

var (
	_ caddy.Provisioner     = (*Subscribe)(nil)
	_ caddyfile.Unmarshaler = (*Subscribe)(nil)
)
