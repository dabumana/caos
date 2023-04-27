FROM golang:1.20

WORKDIR /usr/local/bin
ENV APP_HOME /usr/local/bin/caos

RUN mkdir -p ${APP_HOME}
RUN git clone https://github.com/dabumana/caos ${APP_HOME}
RUN make -C ${APP_HOME}/src build

RUN mv ${APP_HOME}/src/bin/* /usr/local/bin
RUN rm -rf ${APP_HOME}

ENTRYPOINT [ "caos-v.0.0.0" ]
