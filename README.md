# Caddy-NATS
> :warning: **This is still very much a work in progress**

`caddy-nats` is a caddy module that allows the caddy server to interact with a
[nats](https://nats.io/) server. This extension supports multiple patterns:
publish/subscribe, fan in/out, and request reply.

The purpose of this project is to better bridge HTTP based services with NATS
in a pragmatic and straightforward way. If you've been wanting to use NATS, but
have some use cases that still need to use HTTP, this may be a really good
option for you.

## Installation

To use `caddy-nats`, simply run the [xcaddy](https://github.com/caddyserver/xcaddy) build tool to create a `caddy-nats` compatible caddy server.

```sh
xcaddy build --with github.com/codegangsta/caddy-nats
```

## Getting Started

Getting up and running with `caddy-nats` is pretty simple:

First [install NATS](https://docs.nats.io/running-a-nats-service/introduction/installation) and make sure the NATS server is running:

```sh
nats-server
```

Then create your Caddyfile:

```Caddyfile
{
  nats {
    reply hello GET https://localhost/hello
  }
}

localhost {
  route /hello {
    respond "Hello, world"
  }
}
```

After running your custom built `caddy` server with this Caddyfile, you can
then run the `nats` cli tool to test it out:

```sh
nats req hello ""
```

## Connecting to NATS

To connect to `nats`, simply use the `nats` global option in your Caddyfile:

```Caddyfile
{
  nats
}
```

This will connect to the currently selected `nats` context that the `nats` CLI
uses. If there is no `nats` context available, it will connect to the default
`nats` server URL.

You can view your `nats` context information via the `nats` cli:

```sh
nats context list
```

We recommend connecting to a specific context:

```Caddyfile
{
  nats <context>
}
```

## Subscribing to a NATS subject

`caddy-nats` supports subscribing to NATS subjects in a few different flavors,
depending on your needs.

### Subscribe Placeholders

All `subscribe` based directives (`subscribe`, `reply`, `queue_subscribe`, `queue_reply`) support the following `caddy` placeholders in the `method` and `url` arguments:

- `{nats.path}`: The subject of this message, with dots "." replaced with a slash "/" to make it easy to map to a URL.
- `{nats.path.*}`: You can also select a segment of a path ex: `{nats.path.0}` or a subslice of a path: `{nats.path.2:}` or `{nats.path.2:7}` to have ultimate flexibility in what to forward onto the URL path.


### `subscribe`

#### Syntax
```Caddyfile
{
  nats {
    subscribe <subject> <method> <url>
  }
}

```

`subscribe` will subscribe to the specific NATS subject (wildcards are
supported) and forward the NATS payload to the specified URL inside the caddy
web server. This directive does not care about the HTTP response, and is
generally fire and forget.

#### Example

Subscribe to an event stream in NATS and call an HTTP endpoint:

```Caddyfile
subscribe events.> POST https://localhost/nats_events/{nats.path.1:}
```

### `reply`

#### Syntax
```Caddyfile
{
  nats {
    reply <subject> <method> <url>
  }
}
```

`reply` will subscribe to the specific NATS subject and forward the NATS
payload to the specified URL inside the caddy web server. This directive will
then respond back to the nats message with the response body of the HTTP
request.

#### Example

Respond to the `hello.world` NATS subject with the response of the `/hello/world` endpoint.

```Caddyfile
reply hello.world GET https://localhost/hello/world
```

### `queue_subscribe`

#### Syntax
```Caddyfile
{
  nats {
    queue_subscribe <subject> <queue> <method> <url>
  }
}

```
`queue_subscribe` operates the same way as `subscribe`, but subscribes under a NATS [queue group](https://docs.nats.io/nats-concepts/core-nats/queue)

#### Example

Subscribe to a worker queue:

```Caddyfile
queue_subscribe jobs.* workers_queue POST https://localhost/{nats.path}
```

### `queue_reply`

#### Syntax
```Caddyfile
{
  nats {
    queue_subscribe <subject> <queue> <method> <url>
  }
}

```
`queue_reply` operates the same way as `reply`, but subscribes under a NATS [queue group](https://docs.nats.io/nats-concepts/core-nats/queue)

#### Example

Subscribe to a worker queue, and respond to the NATS message:

```Caddyfile
queue_reply jobs.* workers_queue POST https://localhost/{nats.path}
```


## Publishing to a NATS subject

`caddy-nats` also supports publishing to NATS subjects when an HTTP call is
matched within `caddy`, this makes for some very powerful bidirectional
patterns.

### Publish Placeholders

All `publish` based directives (`publish`, `request`) support the following `caddy` placeholders in the `subject` arguments:

- `{nats.subject}`: The path of the http request, with slashes "." replaced with dots "." to make it easy to map to a NATS subject.
- `{nats.subject.*}`: You can also select a segment of a subject ex: `{nats.subject.0}` or a subslice of the subject: `{nats.subject.2:}` or `{nats.subject.2:7}` to have ultimate flexibility in what to forward onto the URL path.

Additionally, since `publish` based directives are caddy http handlers, you also get access to all [caddy http placeholders](https://caddyserver.com/docs/modules/http#docs).

### `nats_publish`

#### Syntax
```Caddyfile
nats_publish [<matcher>] <subject> {
  timeout <timeout-ms>
}
```
`nats_publish` publishes the request body to the specified NATS subject. This
http handler is not a terminal handler, which means it can be used as
middleware (Think logging and events for specific http requests).

#### Example

Publish an event before responding to the http request:

```Caddyfile
localhost {
  route /hello {
    nats_publish events.hello
    respond "Hello, world"
  }
}
```

### `nats_request`

#### Syntax
```Caddyfile
nats_request [<matcher>] <subject> {
  timeout <timeout-ms>
}
```
`nats_request` publishes the request body to the specified NATS subject, and
writes the response of the NATS reply to the http response body.

#### Example

Publish an event before responding to the http request:

```Caddyfile
localhost {
  route /hello/* {
    nats_request hello_service.{nats.subject.1}
  }
}
```

## What's Next?
While this is currently functional and useful as is, here are the things I'd like to add next:

- [ ] Add more examples in the /examples directory
- [ ] Add Validation for all caddy modules in this package
- [ ] Add godoc comments
- [ ] Support mapping nats headers and http headers for upstream and downstream
- [ ] Customizable error handling
- [ ] Jetstream support (for all that persistence babyyyy)
- [ ] nats KV for storage
