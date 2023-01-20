FROM ubuntu:latest
RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y --no-install-recommends golang \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
COPY ./bin/caos /usr/local/bin
ENTRYPOINT ["caos"]