FROM golang:1.20.1-alpine3.17

RUN apk update

RUN apk add --no-cache curl

WORKDIR /go/src/app

COPY . .

RUN go mod tidy \
    && go mod verify

RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

ENTRYPOINT ["air"]