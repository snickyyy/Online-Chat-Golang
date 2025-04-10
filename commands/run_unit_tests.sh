#!/bin/sh

go mod download

cd src/tests/unit && go test -v
