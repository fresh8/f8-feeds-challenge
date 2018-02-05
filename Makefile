build-feed:
	go build -o bin/feed github.com/fresh8/f8-feeds-challenge/feed

run-feed:
	make build-feed && ./bin/feed

build-importer:
	go build -o bin/importer github.com/Roverr/f8-feeds-challenge/importer

run-importer:
	make build-importer && ./bin/importer

test-importer:
	go test ./importer
