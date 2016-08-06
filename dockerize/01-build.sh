#!/bin/sh

DOCKER="docker"
TMP_D="tmp"

BUILD_MACHINE_GOPATH="/go"

BUILD_MACHINE_ID=$($DOCKER run -e GOBIN=${BUILD_MACHINE_GOPATH}/bin -d golang:1.6.3 sleep infinity)

$DOCKER cp "$GOPATH/src" "${BUILD_MACHINE_ID}:${BUILD_MACHINE_GOPATH}/"
$DOCKER exec -i $BUILD_MACHINE_ID go install "$BUILD_MACHINE_GOPATH/src/github.com/rgafiyatullin/imc/server/imcd.go"
$DOCKER cp "${BUILD_MACHINE_ID}:${BUILD_MACHINE_GOPATH}/bin/imcd" "${TMP_D}/"

$DOCKER stop $BUILD_MACHINE_ID
$DOCKER rm $BUILD_MACHINE_ID