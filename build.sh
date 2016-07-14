#!/bin/bash

go get -u -v github.com/kardianos/osext
go get -u -v github.com/russross/blackfriday

rm -rf bin/*
go build -o bin/mdv main.go
cp -r static bin/
