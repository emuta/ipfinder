FROM golang:alpine AS build

WORKDIR /go/src

COPY ./ /go/src

RUN apk update && apk add git \
    && go get github.com/ipipdotnet/datx-go \
    && CGO_ENABLED=0 go build -ldflags '-w -s' -o ipfinder \
    && apk del git && rm -rf /var/cache/apk

FROM scratch

COPY ./ /data
COPY --from=build /go/src/ipfinder /usr/local/bin/ipfinder

CMD ["ipfinder", "-dbfile", "/data/17monipdb.datx"]