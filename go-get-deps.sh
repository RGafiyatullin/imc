#!/bin/sh

export GOPATH="$(pwd)"
go get 'github.com/cyberdelia/go-metrics-graphite' && \
go get 'github.com/rcrowley/go-metrics' && \
go get 'github.com/steveyen/gtreap' && \
go get 'github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre'
