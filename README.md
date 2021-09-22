# zipline

## Ensure GO environment is setup using GOPATH

Even though since 1.11 GOPATH isn't strictly required in favor of modules, but it is required and used by `zipline`. You can read more about this here: https://insujang.github.io/2020-04-04/go-modules/

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
Project must be configured with GOPATH.

`zipline ./...` from your project root

