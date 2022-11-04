#!/bin/sh
go build -o a.exe play/play.go && ./a.exe
rm a.exe