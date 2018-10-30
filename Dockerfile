# Build server
FROM golang:alpine as builder

WORKDIR /vue-12factor
COPY main.go .
RUN go build -ldflags "-w -s"

# Runtime container
#FROM scratch
FROM alpine:latest
WORKDIR /srv/http

COPY --from=builder /vue-12factor/vue-12factor /bin/

EXPOSE 80
ENTRYPOINT ["/bin/vue-12factor"]
#CMD ["/bin/vue-12factor"]

