FROM golang:1.13.8

ENV APP_NAME myproject
ENV PORT 443

COPY . /go/src/${APP_NAME}
WORKDIR /go/src/${APP_NAME}

RUN go get ./
RUN go build -o ${APP_NAME}

CMD ./${APP_NAME}

EXPOSE ${PORT}



#FROM golang:1.11.5

#ENV PORT 80

#RUN mkdir -p /go/src/github.com/MarErm27/webTrainerWithUAdmin
#ADD . /go/src/github.com/MarErm27/webTrainerWithUAdmin

#RUN go get -u github.com/uadmin/uadmin/...

#...

#RUN go install github.com/MarErm27/webTrainerWithUAdmin

#WORKDIR /go/src/github.com/MarErm27/webTrainerWithUAdmin

#RUN go build -o uadmin

#ENTRYPOINT /go/bin

#CMD ./

#EXPOSE ${PORT}