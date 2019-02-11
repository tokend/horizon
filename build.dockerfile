# horizonbuild
FROM golang:1.10

WORKDIR /go/src/gitlab.com/tokend/horizon
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X horizon.version `git rev-parse HEAD`" -o /usr/local/bin/horizon gitlab.com/tokend/horizon/cmd/horizon