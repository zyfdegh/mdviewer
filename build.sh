#!/bin/bash

rm -rf bin/*
go build -o bin/mdv src/main.go
cp -r src/static bin/
cp README.md bin/README.md
