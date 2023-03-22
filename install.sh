#!/bin/zsh
cd $(dirname $(realpath $0)) || return

go build -o ./bin/oncorhynchus ./cmd/oncorhynchus || return

echo 'compiled'

unalias cp 2>/dev/null
mkdir ~/bin 2>/dev/null
cp -rf ./bin/* ~/bin || return

echo 'installed'