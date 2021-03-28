#
# Makefile
#
VERSION = snapshot
GHRFLAGS =
.PHONY: build release

default: build

build:
	# gox -osarch="darwin/amd64 darwin/arm64" -output="pkg/$(VERSION)/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}" ./cmd/...
	gox -osarch="darwin/amd64 darwin/arm64 linux/amd64" -output="pkg/$(VERSION)/{{.OS}}_{{.Arch}}/{{.Dir}}" ./cmd/...

release:
	ghr  -u bilal-bhatti $(GHRFLAGS) v$(VERSION) pkg/$(VERSION)

