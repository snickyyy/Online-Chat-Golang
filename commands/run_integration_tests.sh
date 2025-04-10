#!/bin/sh

go mod download

cd /OCG/src/tests/integration && go test -v

