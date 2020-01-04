FROM golang:1.13.0-stretch AS builder

WORKDIR /build

#COPY go.mod .
#COPY go.sum .
#RUN go mod download

COPY /cmd .

RUN go get github.com/go-telegram-bot-api/telegram-bot-api
RUN go test
RUN go build -o surfbot

WORKDIR /dist
RUN cp /build/surfbot ./surfbot

RUN ldd surfbot | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname ./%); cp % ./%;'
RUN mkdir -p lib64 && cp /lib64/ld-linux-x86-64.so.2 lib64/

RUN mkdir /data

FROM scratch

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --chown=0:0 --from=builder /dist /
COPY --chown=65534:0 --from=builder /data /data
USER 65534
WORKDIR /data

ENTRYPOINT ["/surfbot"]