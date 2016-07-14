#!/bin/bash

rm -rf bin/*
go build -o bin/mdv main.go
cp -r static bin/
