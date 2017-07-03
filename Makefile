build-feed:
	go build -o bin/feed github.com/connectedventures/f8-feeds-challenge/feed

run-feed:
	make build-feed && ./bin/feed
