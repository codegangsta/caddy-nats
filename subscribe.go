package caddynats

import (
	"bytes"
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
	ctx     caddy.Context
	logger  *zap.Logger
	httpApp *caddyhttp.App
}

func init() {
	caddy.RegisterModule(Subscribe{})
}

func (Subscribe) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "nats.handlers.subscribe",
		New: func() caddy.Module { return new(Subscribe) },
	}
}

func (s *Subscribe) Provision(ctx caddy.Context) error {
	s.ctx = ctx
	s.logger = ctx.Logger(s)

	return nil
}

func (s *Subscribe) Subscribe(conn *nats.Conn) error {
	s.logger.Info("Subscribing to NATS subject", zap.String("subject", s.Subject))
	httpAppIface, err := s.ctx.App("http")
	if err != nil {
		return err
	}
	s.httpApp = httpAppIface.(*caddyhttp.App)

	sub, err := conn.Subscribe(s.Subject, s.handler)
	s.sub = sub

	return err
}

func (s *Subscribe) Unsubscribe(conn *nats.Conn) error {
	s.logger.Info("Unsubscribing from NATS subject", zap.String("subject", s.Subject))

	return s.sub.Drain()
}

func (s *Subscribe) handler(msg *nats.Msg) {
	s.logger.Debug("Handling message NATS on subject", zap.String("subject", msg.Subject))

	// TODO: Support multiple servers, for now just pick the first one
	server := s.httpApp.Servers["srv0"]

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
	_ Handler               = (*Subscribe)(nil)
)
