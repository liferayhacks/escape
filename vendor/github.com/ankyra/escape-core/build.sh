#!/bin/bash -e

set -euf -o pipefail

rm -rf docs/generated/
docker rm src || true
docker create -v /go/src/github.com/ankyra/ --name src golang:1.9.0 /bin/true
docker cp "$PWD" src:/go/src/github.com/ankyra/tmp
docker run --rm --volumes-from src \
    -w /go/src/github.com/ankyra/ \
    golang:1.9.0 mv tmp escape-core
docker run --rm \
    --volumes-from src \
    -w /go/src/github.com/ankyra/escape-core \
    golang:1.9.0 bash -c "go build && go run docs/generate_stdlib_docs.go && go run docs/generate_pages.go"
docker cp 'src:/go/src/github.com/ankyra/escape-core/docs/generated' docs/

docker rm src
