# Gouter

**EXPERIMENTAL** - only for development purpose.

Simple development reverse proxy.

You have to install [**Go**](https://go.dev/dl/) to run this package.

For development uses [reflex](https://github.com/cespare/reflex) package.

Reverse proxy running on port :80.

Key commands are written in Makefile, which is recommended to use.

## Config file info
- gouter.yaml must be in the same folder
- hosts use `::` separator

## Config file example
```yaml
hosts:
  - 'http://localhost:8000/test/::/test/'
  - 'http://localhost:9000/::/'
```

## Goals
- [ x ] Working reverse proxy
- [  ] Error logs
- [  ] TLS
- [  ] Load balancing
- [  ] Tests