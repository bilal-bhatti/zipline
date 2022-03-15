# zipline

[![Go Report Card](https://goreportcard.com/badge/github.com/bilal-bhatti/zipline)](https://goreportcard.com/report/github.com/bilal-bhatti/zipline)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/bilal-bhatti/zipline/blob/master/LICENSE)
[![release](https://img.shields.io/github/release/bilal-bhatti/zipline/all.svg)](https://github.com/bilal-bhatti/zipline/releases)
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/https://pkg.go.dev/github.com/bilal-bhatti/zipline)](https://pkg.go.dev/github.com/bilal-bhatti/zipline)

Zipline is a tool to help build RESTful APIs in Go. The intent is to enable separation of HTTP/S request/response marshaling and un-marshaling from request/response processing, apply consistent request/response handling, and generate API documentation, without requiring a runtime dependency on any additional packages. It can be used with any of the popular routing packages. The examples in this repo use `chi`.

It's a template driven code generation tool. An example of a template is [here](https://github.com/bilal-bhatti/zipline/blob/master/example/web/bindings.go), without it the tool doesn't do anything. Once a similar template has been created, simply run `zipline ./...` from the Go project root.

It will generate the following files:
* `bindings_gen.go` [example](https://github.com/bilal-bhatti/zipline/blob/master/example/web/bindings_gen.go)
* `API.md` [example](https://github.com/bilal-bhatti/zipline/blob/master/API.md)
* `api.oasv2.json` [example](https://github.com/bilal-bhatti/zipline/blob/master/api.oasv2.json)
* `api.oasv3.json` [example](https://github.com/bilal-bhatti/zipline/blob/master/api.oasv3.json)

Check [here](https://github.com/bilal-bhatti/zipline/tree/master/example/web) for examples.

## Ensure GO environment is setup properly using `GOPATH` if `GO111MODULE=off`

zipline will attempt to detect package with `bindings.go` file and generate the `bindings_gen.go` file in the same location.

You can read more about the differences and implications here: https://insujang.github.io/2020-04-04/go-modules/

* `export GOPATH=<go directory>`
* `export GOBIN=GOPATH/bin`
* `export PATH=$GOBIN:$PATH`

## Installation
* Brew tap
    * `brew tap bilal-bhatti/homebrew-taps`
    * `brew install zipline`

OR install from source

* Get and install
    * `go install github.com/bilal-bhatti/zipline/cmd/zipline@latest`

## Usage
To explore the examples included in this repo navigate to `cd <path to project>/zipline` and run `zipline ./...` from your project root.
Use `zipline --debug ./...` to enable trace logging to help troubleshoot problems.

## Troubleshooting
VS Code uses `gopls` as the language server and can run into trouble parsing the `bidnings.go` file because of the `ziplinegen` build tag. To fix this, `gopls` needs to be made aware of the build tag, so it can analyze the code properly. Add the following setting to VS Code.

```
“gopls”: {
    “buildFlags”: [“-tags=ziplinegen”],
}
```
