build-cron:
	go build -o bin/cron cmd/cron/main.go

build-server:
	go build -o bin/server cmd/server/main.go

build: build-cron build-server

clean:
	rm -f bin/*

all: build