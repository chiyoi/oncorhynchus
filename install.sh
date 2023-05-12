#!/bin/zsh
cd $(dirname $(realpath $0)) || return

GOBIN=~/bin go install ./cmd/trinity || return

echo 'programs installed'

unalias cp 2>/dev/null
mkdir ~/bin 2>/dev/null
cp -rf ./bin/* ~/bin || return

echo 'scripts installed'
