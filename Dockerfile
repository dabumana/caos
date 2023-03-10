FROM ubuntu:22.04
RUN apt-get update \
    && apt-get install -y --no-install-recommends golang=1.19.5 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
COPY ./bin/caos /usr/local/bin
ENTRYPOINT ["caos"]