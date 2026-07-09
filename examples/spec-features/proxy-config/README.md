# Proxy Configuration

Demonstrates how to configure HTTP proxy for all API requests, including
OAuth2 token exchange and Digest challenge requests.

## Supported proxy schemes

| Scheme | Description | DNS resolution |
|--------|-------------|----------------|
| `http://` | HTTP proxy | Local |
| `https://` | HTTPS proxy | Local |
| `socks5://` | SOCKS5 proxy | Local |
| `socks5h://` | SOCKS5 proxy | **Through proxy** |

## What it demonstrates

- `proxy.url` — proxy server URL
- `proxy.username` / `proxy.password` — proxy authentication (optional)
- `proxy.bypass` — list of hosts/networks to exclude from proxying
- `random: true` — enables browser-like headers (User-Agent, Accept, Referer, Sec-*)

## Bypass patterns

- `*.local` — wildcard: matches any subdomain of `.local`
- `10.0.0.0/8` — CIDR notation for private networks
- `example.com` — exact host match

## Expected behavior

- All external API requests go through the SOCKS5 proxy
- Internal/local requests bypass the proxy
- OAuth2 token requests also use the proxy
- Browser-like headers are randomly generated at startup
