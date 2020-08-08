.DEFAULT_GOAL := all

build:
	go build -o main .
	rm -rf dist
	mkdir dist
	cp main dist/
	cp assets/* dist

install:
	rm -rf "$(HOME)/Library/Application Support/com.elgato.StreamDeck/Plugins/us.hervalicio.gameoflife.sdPlugin"
	mv dist "$(HOME)/Library/Application Support/com.elgato.StreamDeck/Plugins/us.hervalicio.gameoflife.sdPlugin"

all: build install;
