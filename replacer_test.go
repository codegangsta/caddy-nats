package caddynats

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"

	"github.com/caddyserver/caddy/v2"
	"github.com/nats-io/nats.go"
)

// Returns a basic request for testing
func reqPath(path string) *http.Request {
	req, _ := http.NewRequest("GET", "http://localhost"+path, &bytes.Buffer{})
	return req
}

func TestAddNatsPublishVarsToReplacer(t *testing.T) {
	type test struct {
		req *http.Request

		input string
		want  string
	}

	tests := []test{

		// Basic subject mapping
		{req: reqPath("/foo/bar"), input: "{nats.subject}", want: "foo.bar"},
		{req: reqPath("/foo/bar/bat/baz"), input: "{nats.subject}", want: "foo.bar.bat.baz"},
		{req: reqPath("/foo/bar/bat/baz?query=true"), input: "{nats.subject}", want: "foo.bar.bat.baz"},
		{req: reqPath("/foo/bar"), input: "prefix.{nats.subject}.suffix", want: "prefix.foo.bar.suffix"},

		// Segment placeholders
		{req: reqPath("/foo/bar"), input: "{nats.subject.0}", want: "foo"},
		{req: reqPath("/foo/bar"), input: "{nats.subject.1}", want: "bar"},

		// Segment Ranges
		{req: reqPath("/foo/bar/bat/baz"), input: "{nats.subject.0:}", want: "foo.bar.bat.baz"},
		{req: reqPath("/foo/bar/bat/baz"), input: "{nats.subject.1:}", want: "bar.bat.baz"},
		{req: reqPath("/foo/bar/bat/baz"), input: "{nats.subject.2:}", want: "bat.baz"},
		{req: reqPath("/foo/bar/bat/baz"), input: "{nats.subject.1:3}", want: "bar.bat"},
		{req: reqPath("/foo/bar/bat/baz"), input: "{nats.subject.0:3}", want: "foo.bar.bat"},
		{req: reqPath("/foo/bar/bat/baz"), input: "{nats.subject.0:4}", want: "foo.bar.bat.baz"},
		{req: reqPath("/foo/bar/bat/baz"), input: "{nats.subject.:3}", want: "foo.bar.bat"},

		// Out of bounds ranges
		{req: reqPath("/foo/bar/bat/baz"), input: "{nats.subject.0:18}", want: ""},
		{req: reqPath("/foo/bar/bat/baz"), input: "{nats.subject.-1:}", want: ""},
	}

	for _, tc := range tests {
		repl := caddy.NewReplacer()
		addNATSPublishVarsToReplacer(repl, tc.req)
		got := repl.ReplaceAll(tc.input, "")
		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("expected: %v, got: %v", tc.want, got)
		}
	}
}

func TestAddNatsSubscribeVarsToReplacer(t *testing.T) {
	type test struct {
		msg *nats.Msg

		input string
		want  string
	}

	tests := []test{
		// Basic subject mapping
		{msg: nats.NewMsg("foo.bar"), input: "{nats.path}", want: "foo/bar"},
		{msg: nats.NewMsg("foo.bar.bat.baz"), input: "{nats.path}", want: "foo/bar/bat/baz"},
		{msg: nats.NewMsg("foo.bar"), input: "prefix/{nats.path}/suffix", want: "prefix/foo/bar/suffix"},

		// // Segment placeholders
		{msg: nats.NewMsg("foo.bar"), input: "{nats.path.0}", want: "foo"},
		{msg: nats.NewMsg("foo.bar"), input: "{nats.path.1}", want: "bar"},

		// // Segment Ranges
		{msg: nats.NewMsg("foo.bar.bat.baz"), input: "{nats.path.0:}", want: "foo/bar/bat/baz"},
		{msg: nats.NewMsg("foo.bar.bat.baz"), input: "{nats.path.1:}", want: "bar/bat/baz"},
		{msg: nats.NewMsg("foo.bar.bat.baz"), input: "{nats.path.2:}", want: "bat/baz"},
		{msg: nats.NewMsg("foo.bar.bat.baz"), input: "{nats.path.1:3}", want: "bar/bat"},
		{msg: nats.NewMsg("foo.bar.bat.baz"), input: "{nats.path.0:3}", want: "foo/bar/bat"},
		{msg: nats.NewMsg("foo.bar.bat.baz"), input: "{nats.path.:3}", want: "foo/bar/bat"},

		// Out of bounds ranges
		{msg: nats.NewMsg("foo.bar.bat.baz"), input: "{nats.path.0:18}", want: ""},
		{msg: nats.NewMsg("foo.bar.bat.baz"), input: "{nats.path.:18}", want: ""},
		{msg: nats.NewMsg("foo.bar.bat.baz"), input: "{nats.path.-1:}", want: ""},
	}

	for _, tc := range tests {
		repl := caddy.NewReplacer()
		addNatsSubscribeVarsToReplacer(repl, tc.msg)
		got := repl.ReplaceAll(tc.input, "")
		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("expected: %v, got: %v", tc.want, got)
		}
	}
}
