FROM alpine:latest
# Define local binary folder
ENV APP_HOME /usr/local/bin/caos
# Update local package manager and install required dependencies
RUN apk update && apk upgrade
RUN apk add git musl-dev make pkgconfig gcc go
RUN apk add curl=7.79.1-r1 curl-dev=7.79.1-r1 --repository=http://dl-cdn.alpinelinux.org/alpine/v3.12/main
# Create source folder and clone the repository
RUN mkdir -p ${APP_HOME}
RUN git clone https://github.com/dabumana/caos ${APP_HOME}/source
WORKDIR ${APP_HOME}/source
# Build the application and delete the source folder
RUN make build
RUN cp ${APP_HOME}/source/src/bin/caos/caos /bin
RUN rm -rf ${APP_HOME}
# Setup argument variables
ARG KEY
ARG ZERO-KEY
# Initialize environment variables
ENV API_KEY $KEY
ENV ZERO_API_KEY $ZERO-KEY
# Define entrypoint
ENTRYPOINT [ "caos" ]
