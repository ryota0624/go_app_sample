run:
	go run server.go

watch-run:
	godo --rebuild && godo server -w
