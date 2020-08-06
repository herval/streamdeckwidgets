.DEFAULT_GOAL := all

logs:
	tail -f ~/Library/Logs/StreamDeck/*

build:
	go build -o main .
	rm -rf dist
	mkdir dist
	cp main dist/
	cp assets/* dist

install:
	rm -rf "$(HOME)/Library/Application Support/com.elgato.StreamDeck/Plugins/us.hervalicio.miccontrol.sdPlugin"
	mv dist "$(HOME)/Library/Application Support/com.elgato.StreamDeck/Plugins/us.hervalicio.miccontrol.sdPlugin"

all: build install;
