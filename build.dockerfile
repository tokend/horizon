# horizonbuild
FROM golang:1.16.2-stretch

ARG version="dirty"

WORKDIR /go/src/gitlab.com/tokend/horizon
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=${version}" -o /usr/local/bin/horizon gitlab.com/tokend/horizon/cmd/horizon 
