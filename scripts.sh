#!/bin/sh
cd $(dirname $(realpath $0))
usage() {
    pwd
    echo "scripts:"
    echo "$0 install"
    echo "    Install the commands."
    echo "$0 path"
    echo "    Show path to the directory where the commands installed."
}

install() {
	test -d ~/bin || mkdir ~/bin || return
	cp -f commands/*  ~/bin || return
	echo "Scripts installed~"
	echo "Add to PATH:"
    path
}

path() {
    echo "$HOME/bin"
}

case $1 in
""|-h|-help|--help)
usage
exit
;;
install|path)
;;
*)
usage
exit 1
esac

$@
