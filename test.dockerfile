FROM golang:1.9

WORKDIR /go/src/gitlab.com/swarmfund/horizon
COPY . .
ENTRYPOINT ["go", "test"]