# GitHub API Proxy to NATS

The Caddyfile in this directory is an example of how to run a reverse proxy to the GitHub API and expose it as service in NATS.

## Testing out the GitHub API in NATS

1. Build a caddy server with `caddy-nats installed`:
```sh
xcaddy build --with github.com/codegangsta/caddy-nats
```

2. Generate a [personal access token](https://github.com/settings/tokens) From GitHub

3. Set the `GITHUB_AUTH_HEADER_TOKEN` environment variable
```sh
GITHUB_AUTH_HEADER_TOKEN=$(echo "<username>:<access_token>" | base64)
```

4. Run the caddy server you just built
```sh
./caddy run --watch
```

5. Call the GitHub API from nats
```sh
# similar to curl https://api.github.com/users/codegangsta
nats req github.GET.users.codegangsta ""
```
