FROM golang:1.20

ENV APP_HOME /go/src/caos
RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME

RUN git clone https://github.com/dabumana/caos $APP_HOME 
RUN cd $APP_HOME \
    go clean \
    go build

ENTRYPOINT [ $APP_HOME/caos ]