FROM golang:1.16-alpine as golang
WORKDIR /app
COPY . . 

RUN go build

FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates

FROM alpine:3.6
EXPOSE 3000

ENV GODEBUG netdns=go

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=golang /app/drone-webhook /bin/

ENTRYPOINT ["/bin/drone-webhook"]