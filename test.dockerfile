FROM golang:1.9

WORKDIR /go/src/gitlab.com/tokend/horizon
COPY . .
ENTRYPOINT ["go", "test"]
