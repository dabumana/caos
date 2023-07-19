FROM golang:1.19

WORKDIR /usr/local/bin
ENV APP_HOME /usr/local/bin/caos

RUN mkdir -p ${APP_HOME}
RUN git clone https://github.com/dabumana/caos ${APP_HOME}

RUN cp ${APP_HOME}/src/bin/caos/* /usr/local/bin
RUN rm -rf ${APP_HOME}

ENTRYPOINT [ "caos" ]
