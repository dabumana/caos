FROM golang:1.20

WORKDIR /usr/local/bin/caos

ENV APP_HOME /usr/local/bin/caos

RUN mkdir -p ${APP_HOME}
RUN git clone https://github.com/dabumana/caos ${APP_HOME}
RUN make -C ${APP_HOME}/src build

RUN mv ${APP_HOME}/src/bin/* ${APP_HOME}
RUN rm -rf ${APP_HOME}/!"caos"

CMD [ ".\caos" ]
