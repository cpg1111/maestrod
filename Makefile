all: build
get-deps:
	echo;
build:
	glide install
	go build -o maestrod main.go
	go build -buildmode=plugin -o docker.so manager/docker/plugin.go
	go build -buildmode=plugin -o kube.so manager/k8s/plugin.go
test:
	ETCD2_SERVICE_HOST=127.0.0.1 ETCD2_SERVICE_PORT=22379 go test ./datastore/etcd/v2/...
	ETCD3_SERVICE_HOST=127.0.0.1 ETCD3_SERVICE_PORT=32379 go test ./datastore/etcd/v3/...
	MONGO_SERVICE_HOST=`docker-machine ip maestrod-dev` MONGO_SERVICE_PORT=27017 go test ./datastore/mongodb/...
	REDIS_SERVICE_HOST=`docker-machine ip maestrod-dev` REDIS_SERVICE_PORT=6379 go test ./datastore/redis/...
	go test ./config/...
	go test ./lifecycle/...
	# go test ./manager/docker/...
	TEST_K8S_URL=https://`dev_deps/minikube ip`:8443 \
	TEST_MAESTRO_VER=latest \
	TEST_BRANCH=master \
	TEST_CONF_PATH=./example.conf.toml \
	TEST_PREV_COMMIT=f5c1e92536b56b09b7cca764a066a1fc3f19cc8d \
	TEST_CURR_COMMIT=02eeac380bae358c6c4f19e720d42e6822cb4903 \
	TEST_CLIENT_CERT=~/.minikube/apiserver.crt \
	TEST_CLIENT_KEY=~/.minikube/apiserver.key \
	ROOT_CA_PATH=~/.minikube/ca.crt \
	go test ./manager/k8s/...
install:
	mkdir -p /opt/bin/maestrod
	mkdir -p /etc/maestrod/
	cp maestrod /opt/bin/maestrod/maestrod
	cp example.conf.toml /etc/maestrod/conf.toml
	mkdir -p /etc/maestrod/conf.d/ /etc/maestrod/plugin.d/
	cp docker.so /etc/maestrod/plugin.d/
	cp kube.so /etc/maestrod/plugin.d/
docker:
	docker build -t maestrod-build -f Dockerfile_build .
	docker run -v `pwd`/dist/:/opt/bin/maestrod/ maestrod-build
	docker build -t maestrod .
