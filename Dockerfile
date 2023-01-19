FROM ubuntu
RUN apt-get update \
    && sudo apt-get upgrade -y \
    && apt-get install golang
ADD ./bin/caos /usr/local/bin
ENTRYPOINT ["caos"]