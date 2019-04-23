FROM golang:alpine AS build

WORKDIR /go/src

COPY ./ /go/src

RUN apk update && apk add git upx \
    && go get github.com/ipipdotnet/datx-go \
    && go get github.com/sirupsen/logrus \
    && CGO_ENABLED=0 go build -ldflags '-w -s' -o ipfinder \
    && upx ipfinder \
    && apk del git && rm -rf /var/cache/apk

FROM scratch

COPY ./ /data
COPY --from=build /go/src/ipfinder /usr/local/bin/ipfinder

CMD ["ipfinder", "-dbfile", "/data/17monipdb.datx"]