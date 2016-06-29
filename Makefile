all: build
get-deps:
	cd $(GOPATH)
	go get github.com/kardianos/govendor
build:
	govendor sync
	go build -o maestrod-container container/main.go
	go build -o maestrod main.go
test:
	go test ./...
install:
	mkdir -p /opt/bin/maestrod
	mkdir -p /etc/maestrod/
	cp maestrod-container /opt/bin/maestrod/maestrod-container
	cp maestrod /opt/bin/maestrod/maestrod
	cp example.conf.toml /etc/maestrod/conf.toml
