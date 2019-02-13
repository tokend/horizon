# horizonbuild
FROM golang:1.10-stretch

WORKDIR /go/src/gitlab.com/tokend/horizon
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=`[[ -z "$(git tag -l --points-at HEAD)" ]] && git rev-parse HEAD || git tag -l --points-at HEAD`" -o /usr/local/bin/horizon gitlab.com/tokend/horizon/cmd/horizon