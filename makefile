APP=caos
# Enter your credentials for OpenAI && ZeroGPT
KEY="<YOUR-API-KEY>"
ZKEY="<YOUR-API-KEY>"
# Assign resources for service pod
CPU=2
# Configuration path
CONFIG_PATH=./ci/service

build: 
	make -C ./src clean
	make -C ./src build
	make -C ./src test

clean:
	make -C ./src clean

coverage:
	make -C ./src coverage

install: build
	make -C ./src install

run: build
	make -C ./src run

test:
	make -C ./src test

tidy:
	make -C ./src tidy

vendor:
	make -C ./src vendor

build-pod:
	docker build --build-arg KEY=${KEY} --build-arg ZKEY=${ZKEY} --pull --rm -f "ci/service/Dockerfile" -t ${APP}:latest ${CONFIG_PATH}

run-pod: build-pod

	docker run -it --cpus=${CPU} ${APP}:latest
