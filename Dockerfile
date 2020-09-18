FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache make

WORKDIR /go/src/exceldumper/
COPY . .

RUN make server

FROM alpine:latest as worker

WORKDIR /opt/bin/
COPY --from=builder /go/src/exceldumper/ .

EXPOSE 3000/tcp

ENTRYPOINT ["./srv"]
