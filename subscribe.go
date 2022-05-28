package caddynats

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Subscribe struct {
	Subject string `json:"subject,omitempty"`
	Method  string `json:"method,omitempty"`
	URL     string `json:"path,omitempty"`

	WithReply bool `json:"with_reply,omitempty"`

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
	s.logger.Info("subscribing to NATS subject", zap.String("subject", s.Subject))
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
	s.logger.Info("unsubscribing from NATS subject", zap.String("subject", s.Subject))

	return s.sub.Drain()
}

func (s *Subscribe) handler(msg *nats.Msg) {
	s.logger.Debug("handling message NATS on subject", zap.String("subject", msg.Subject))

	req, err := http.NewRequest(s.Method, s.URL, bytes.NewBuffer(msg.Data))
	if err != nil {
		s.logger.Error("error creating request", zap.Error(err))
		return
	}

	server, err := s.matchServer(s.httpApp.Servers, req)
	if err != nil {
		s.logger.Error("error matching server", zap.Error(err))
		return
	}

	if s.WithReply {
		rec := httptest.NewRecorder()
		server.ServeHTTP(rec, req)
		//TODO Handle error
		msg.Respond(rec.Body.Bytes())
		return
	}

	server.ServeHTTP(noopResponseWriter{}, req)
}

func (s *Subscribe) matchServer(servers map[string]*caddyhttp.Server, req *http.Request) (*caddyhttp.Server, error) {
	repl := caddy.NewReplacer()
	for _, server := range servers {
		r := caddyhttp.PrepareRequest(req, repl, nil, server)
		for _, route := range server.Routes {
			if route.MatcherSets.AnyMatch(r) {
				return server, nil
			}
		}
	}

	return nil, fmt.Errorf("no server matched for the current url: %s", req.URL)
}

var (
	_ caddy.Provisioner = (*Subscribe)(nil)
	_ Handler           = (*Subscribe)(nil)
)
