build-cron:
	go build -o bin/survilleray-cron cmd/cron/main.go

build-server:
	go build -o bin/survilleray-server cmd/server/main.go

build-migrate:
	go build -o bin/survilleray-migrate cmd/migrate/main.go

build-job:
	go build -o bin/survilleray-job cmd/survilleray/main.go


build: build-cron build-server build-migrate build-job

clean:
	rm -f bin/*

all: build