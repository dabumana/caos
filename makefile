APP=caos
VERSION=v.0.0.0
# Configuration path
CONFIG_PATH=./ci/service

build:
	make -C ./src clean
	make -C ./src test
	make -C ./src build APP=${APP} VERSION=${VERSION}

test:
	make -C ./src test

clean:
	make -C ./src clean

coverage:
	make -C ./src coverage

run: build
	make -C ./src run

vendor:
	make -C ./src vendor

build-pod:
	docker build --no-cache -t ${APP} ${CONFIG_PATH} 

run-pod: build-pod

	docker run ${APP}
