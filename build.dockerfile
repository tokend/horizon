# horizonbuild
FROM golang:1.9

WORKDIR /go/src/gitlab.com/tokend/horizon
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/horizon gitlab.com/tokend/horizon/cmd/horizon
