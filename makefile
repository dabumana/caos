APP=caos
CONFIG_PATH=./ci/service

build-pod:
	docker build ${CONFIG_PATH} -t caos

deploy: build
	docker run ${APP} --env-file=.env 
