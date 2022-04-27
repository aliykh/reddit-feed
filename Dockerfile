ARG PROXY='https://proxy.golang.org,direct'
FROM golang:1.16-alpine as build

WORKDIR /app

RUN export GOPROXY=$PROXY

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o bin/reddit ./cmd

FROM alpine:latest

WORKDIR /root/

RUN mkdir config
COPY --from=build /app/config/ ./config/
COPY --from=build /app/bin/reddit .

CMD [ "./reddit" ]