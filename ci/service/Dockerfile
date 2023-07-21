FROM alpine:latest

ENV APP_HOME /usr/local/bin/caos

RUN apk update && apk upgrade
RUN apk add git musl-dev make pkgconfig gcc go
RUN apk add curl=7.79.1-r1 curl-dev=7.79.1-r1 --repository=http://dl-cdn.alpinelinux.org/alpine/v3.12/main

RUN mkdir -p ${APP_HOME}
RUN git clone https://github.com/dabumana/caos ${APP_HOME}/source
WORKDIR ${APP_HOME}/source

RUN make build
RUN cp ${APP_HOME}/source/src/bin/caos/caos /bin
RUN rm -rf ${APP_HOME}

ENTRYPOINT [ "caos" ]
