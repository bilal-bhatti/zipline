# Release using GoReleaser

## Create release tag
`git tag -a v0.1.0 -m "Release title"`

## Push tag to git
`git push origin v0.1.0`

## Install and run GoReleaser
``` sh
brew install goreleaser
goreleaser check
goreleaser release --clean
```

## reset tags
``` sh
git tag -d v0.1.0
git push origin :v0.1.0
```

## Copy Pasta
``` sh
git tag -a v0.5.0 -m ""
git push origin v0.5.0
goreleaser release --clean
```
