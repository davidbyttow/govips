#!/usr/bin/env bash
set -x

git config --global url."git@github.com:".insteadOf https://github.com/
ssh-keyscan github.com >> ~/.ssh/known_hosts

GOPATH=`pwd`/deps go mod download

