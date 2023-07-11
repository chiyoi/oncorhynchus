.PHONY = install

all:
	@echo "No object to build~"

install:
	test -d ~/bin || mkdir ~/bin
	cp -f cmd/*  ~/bin
	@echo
	@echo "Scripts installed~"
	@echo "Add to PATH: $$HOME/bin"
