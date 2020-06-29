#!/bin/sh

mkdir -p /usr/share/icons/secprac
cp img/* /usr/share/icons/secprac

# will download latest release instead of compile when we have one
go install
