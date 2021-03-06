FROM golang:1.16-alpine as build

ARG PROXY=https://proxy.golang.org,direct

WORKDIR /app

ENV GOPROXY=${PROXY}

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