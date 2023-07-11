#!/usr/bin/env sh
usage() {
    echo "scripts:"
    echo "scripts.sh install"
    echo "    Install the commands."
    echo "scripts.sh path"
    echo "    Show path to the directory where the commands installed."
}

install() {
	test -d ~/bin || mkdir ~/bin
	cp -f cmd/*  ~/bin
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
