APP=caos
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

run: build
	make -C ./src run

test:
	make -C ./src test

tidy:
	make -C ./src tidy

vendor:
	make -C ./src vendor

build-pod:
	docker build --no-cache -t ${APP} ${CONFIG_PATH} 

run-pod: build-pod

	docker run ${APP}
