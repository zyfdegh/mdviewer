#!/bin/bash

go get -u -v github.com/kardianos/osext
go get -u -v github.com/russross/blackfriday
go get -u -v github.com/spf13/cobra

rm -rf bin/*
go build -o bin/mdv main.go
cp -r static bin/
