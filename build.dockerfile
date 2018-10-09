FROM golang:1.9
WORKDIR /go/src/gitlab.com/tokend/horizon
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /binary gitlab.com/tokend/horizon/cmd/horizon
RUN cp ./run.docker /run.docker

FROM ubuntu:latest
COPY --from=0 /binary .
COPY --from=0 /run.docker .
RUN chmod +x /run.docker
ENTRYPOINT ["/run.docker"]
