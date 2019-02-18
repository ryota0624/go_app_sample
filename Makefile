run:
	go run server.go

watch-run:
	godo server -w --rebuild

build-githook-post:
	go build -o bin/githooks/post_commit githooks/post_commit.go

install-githook-post:
	go install  -v ./githooks/post_commit