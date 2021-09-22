# zipline

Zipline is a tool to help build RESTful APIs in Go. The intent is to enable separation of HTTP/S request/response marshaling and un-marshaling from request/response processing, apply consistent request/response handling, and generate API documentation, without requiring a runtime dependency on any additional packages. It can be used with any of the popular routing packages. The examples in this repo use `chi`.

It's a template driven code generation tool. An example of a template is [here](https://github.com/bilal-bhatti/zipline/blob/master/example/web/bindings.go), without it the tool doesn't do anything. Once a similar template has been created, simply run `zipline ./...` from the Go project root.

It will generate the following files:
* `bindings_gen.go` in the same folder as the `bindings.go`
* `API.md` [example](https://github.com/bilal-bhatti/zipline/blob/master/API.md)
* `api.oasv2.json` [example](https://github.com/bilal-bhatti/zipline/blob/master/api.oasv2.json)
* `api.oasv3.json` [example](https://github.com/bilal-bhatti/zipline/blob/master/api.oasv3.json)

Check [here](https://github.com/bilal-bhatti/zipline/tree/master/example/web) for examples.

## Ensure GO environment is setup using GOPATH

Even though since 1.11 GOPATH isn't strictly required in favor of modules, but it is required and used by `zipline`. You can read more about the differences and implications here: https://insujang.github.io/2020-04-04/go-modules/

* `export GOPATH=<go directory>`
* `export GOBIN=GOPATH/bin`
* `export PATH=$GOBIN:$PATH`

## Installation
* Brew tap
    * `brew tap bilal-bhatti/homebrew-zipline`
    * `brew install zipline`

OR install from source

* Get and install
    * `go get -u github.com/bilal-bhatti/zipline/cmd/zipline`
    * `go install github.com/bilal-bhatti/zipline/cmd/zipline`

## Usage
Project must be configured with GOPATH, i.e. project path/hierarchy should look like this `<GOPATH>/src/github.com/<git org/user>/<repo>`. For this repo it would be `<GOPATH>/src/github.com/bilal-bhatti/zipline`.

`zipline ./...` from your project root

