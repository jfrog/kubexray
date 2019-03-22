FROM golang:1.12-alpine AS builder

ARG srcpath="/build/kubexray"

RUN apk --no-cache add git && \
    mkdir -p "$srcpath"

ADD cmd/kubexray/ "$srcpath"

RUN cd "$srcpath" && \
    GO111MODULE=on GOPROXY=https://gocenter.io CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --installsuffix cgo --ldflags="-s" -o /kubexray

FROM alpine:3.9
RUN apk --no-cache add --update ca-certificates

COPY --from=builder /kubexray /bin/kubexray

# Create user
ARG uid=1000
ARG gid=1000
RUN addgroup -g $gid kubexray && \
    adduser -D -u $uid -G kubexray kubexray

USER kubexray

ENTRYPOINT ["/bin/kubexray"]
