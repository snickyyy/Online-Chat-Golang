#!/bin/sh

go mod download

cd src && go run main.go
