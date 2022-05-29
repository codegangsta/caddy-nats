# GitHub API Proxy to NATS

The Caddyfile in this directory is an example of how to run a reverse proxy to the GitHub API and expose it as service in NATS.

To run this simply run a version of `caddy` with the `caddy-nats` extension supported in this directory.

## Testing out the GitHub API in NATS

Once the `caddy` server and `nats-server` is running, you can use the `nats` cli to test out the proxy:

``` sh
# similar to curl https://api.github.com/zen
nats req github.GET.zen ""
```

## Authentication
TODO: Need to add authentication support and a guide on base64 encoding credentials
