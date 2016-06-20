all: build
get-deps:
	cd $(GOPATH)
	go get github.com/tools/godep
	# cd $(GOPATH)src/github.com/cpg1111/
	# wget https://github.com/cpg1111/maestro/archive/master.zip
	# unzip master.zip
	# mv maestro-master $(GOPATH)src/github.com/cpg1111/maestro
	go get -d github.com/cpg1111/maestro
	cd $(GOPATH)/src/github.com/cpg1111/maestro/
	docker build -t maestro_build -f Dockerfile_build .
	docker run --rm -v "$(pwd)":/go/src/github.com/cpg1111/maestro/ maestro_build
	cp ./dist/maestro /usr/bin/maestro
	mkdir /etc/maestro/
	cp ./test_conf.toml /etc/maestro/conf.toml
	cd -
build:
	godep restore ./...
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
