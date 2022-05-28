package caddynats

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

const publishDefaultTimeout = 10000

func init() {
	caddy.RegisterModule(Publish{})
}

type Publish struct {
	Subject   string `json:"subject,omitempty"`
	WithReply bool   `json:"with_reply,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
	Timeout   int64  `json:"timeout,omitempty"`

	logger *zap.Logger
	app    *App
}

func (Publish) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.nats_publish",
		New: func() caddy.Module { return new(Publish) },
	}
}

func (p *Publish) Provision(ctx caddy.Context) error {
	p.logger = ctx.Logger(p)

	natsAppIface, err := ctx.App("nats")
	if err != nil {
		return fmt.Errorf("getting NATS app: %v. Make sure NATS is configured in global options", err)
	}

	p.app = natsAppIface.(*App)

	return nil
}

func (p Publish) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	repl := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)
	prefix := repl.ReplaceAll(p.Prefix, "")
	addNATSVarsToReplacer(repl, r, w, prefix)

	//TODO: What method is best here? ReplaceAll vs ReplaceWithErr?
	subj := repl.ReplaceAll(p.Subject, "")

	//TODO: Check max msg size
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	p.logger.Debug("Publishing NATS message", zap.String("subject", subj), zap.Bool("with_reply", p.WithReply), zap.Int64("timeout", p.Timeout))

	if p.WithReply {
		return p.natsRequestReply(subj, data, w)
	}

	// Otherwise. just publish like normal
	err = p.app.conn.Publish(subj, data)
	if err != nil {
		return err
	}

	return next.ServeHTTP(w, r)
}

func (p Publish) natsRequestReply(subject string, reqBody []byte, w http.ResponseWriter) error {
	//TODO: Configurable timeout
	m, err := p.app.conn.Request(subject, reqBody, time.Duration(p.Timeout)*time.Millisecond)
	if err != nil {
		return err
	}

	// TODO: Make error handlers configurable
	if err == nats.ErrNoResponders {
		w.WriteHeader(http.StatusNotFound)
		return err
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	_, err = w.Write(m.Data)

	return err
}

var (
	_ caddyhttp.MiddlewareHandler = (*Publish)(nil)
	_ caddy.Provisioner           = (*Publish)(nil)
	_ caddyfile.Unmarshaler       = (*Publish)(nil)
)
