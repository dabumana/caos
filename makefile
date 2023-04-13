APP=caos
ARCH=amd64
CONFIG_PATH=./ci/service
VERSION=v.0.2.0

build:
	make -C ./src clean
	make -C ./src test
	make -C ./src build APP=${APP} ARCH=${ARCH} VERSION=${VERSION}

test:
	make -C ./src test

clean:
	make -C ./src clean

coverage:
	make -C ./src coverage

run:
	make -C ./src run

deploy-pod:
	docker build ${CONFIG_PATH} -t caos

run-pod: build
	docker run ${APP} --env-file=.env 
