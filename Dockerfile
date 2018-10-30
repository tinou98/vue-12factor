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

ARG BUILD_DATE
ARG VCS_REF
ARG VERSION
LABEL org.label-schema.build-date=$BUILD_DATE \
	org.label-schema.name="vue-12factor" \
	org.label-schema.description="A container that helps to resolve the run-time configuration for a pre-built Web application." \
	org.label-schema.vcs-ref=$VCS_REF \
	org.label-schema.vcs-url="https://github.com/tinou98/vue-12factor" \
	org.label-schema.vendor="MATILLAT Quentin" \
	org.label-schema.version=$VERSION \
	org.label-schema.schema-version="1.0"
#	org.label-schema.url="website"
