# Release using GoReleaser

## Create release tag
`git tag -a v0.1.0 -m "Release title"`


## Push tag to git
`git push origin v0.1.0`


## Install and run GoReleaser
`brew install goreleaser`

`goreleaser check`

`goreleaser release --rm-dist`

## reset tags
`git tag -d v0.1.0`

`git push origin :v0.1.0`

