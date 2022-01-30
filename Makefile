default:
	go build -ldflags="-s -w" -o ./bin/goldencli ./cmd/goldencli/
	strip ./bin/goldencli

install:
	cp ./bin/goldencli ~/bin/