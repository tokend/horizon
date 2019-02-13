# horizonbuild
FROM golang:1.10-stretch

WORKDIR /go/src/gitlab.com/tokend/horizon
COPY . .
RUN CGO_ENABLED=0 GOOS=linux ./build.sh 
