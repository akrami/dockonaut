FROM golang:alpine

COPY . /go/src/dockonaut
WORKDIR /go/src/dockonaut
RUN go build -o dist/dockonaut cmd/dockonaut/main.go

FROM alpine

RUN apk add --update docker docker-compose git openrc && \
  rc-update add docker boot

VOLUME /var/run/docker.sock

COPY --from=0 /go/src/dockonaut/dist/dockonaut /usr/bin/dockonaut
CMD ["sleep", "infinity"]
