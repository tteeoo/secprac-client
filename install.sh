#!/bin/sh

mkdir -p /usr/local/share/secprac
cp img/* /usr/local/share/secprac

# will download latest release instead of compile when we have one
go install
