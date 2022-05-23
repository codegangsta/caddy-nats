# Caddy-NATS

> :warning: **Still a WIP**

## Example Caddyfile directives
```sh
# connect to nats default server
nats

# Connect to nats server
nats demo.nats.io 

# connect to nats with custom options
nats demo.nats.io {
  user jeremy
  password foobar
}

# Publishing a nats message when a request is made (middleware, does not write a response)
nats_publish [<matcher>] <subject>

# subscribing to a nats subject, does not write a response to nats
nats_subscribe <subject> <method> <path>

# Publishes a request to nats, the response from nats will be written as an http response
nats_request [<matcher>] <subject>

# Replies to a nats subject with the result of the http request to method and path
nats_reply <subject> <method> <path>
```
