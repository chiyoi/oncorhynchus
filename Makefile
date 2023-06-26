.PHONY = install

all:
	@echo 'No object to build~'

install:
	test -d ~/bin || mkdir ~/bin
	cp -f sources/* ~/bin
	@echo 'Scripts installed~'
