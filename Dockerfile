FROM golang:1.11-alpine AS builder
MAINTAINER "<solutions@jfrog.com>"

ARG srcpath="/build/kube-xray"

RUN apk --no-cache add git && \
    mkdir -p "$srcpath"

ADD . "$srcpath"

RUN cd "$srcpath" && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --installsuffix cgo --ldflags="-s" -o /kube-xray

FROM alpine:3.8
RUN apk --no-cache add --update ca-certificates

COPY --from=builder /kube-xray /bin/kube-xray

# Create user
ARG uid=1000
ARG gid=1000
RUN addgroup -g $gid kube-xray && \
    adduser -D -u $uid -G kube-xray kube-xray

USER kube-xray

ENTRYPOINT ["/bin/kube-xray"]
