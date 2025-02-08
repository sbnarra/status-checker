FROM scratch AS status
ARG CREATED
ARG REVISION
# https://github.com/opencontainers/image-spec/blob/main/annotations.md
LABEL org.opencontainers.image.ref.name "sbnarra/status-checker"
LABEL org.opencontainers.image.title "status-checker"
LABEL org.opencontainers.image.description "status checker"
LABEL org.opencontainers.image.source "https://github.com/sbnarra/status-checker"
LABEL org.opencontainers.image.documentation "https://github.com/sbnarra/status-checker"
LABEL org.opencontainers.image.created ${CREATED:-unset}
LABEL org.opencontainers.image.version ${REVISION:-unset}
LABEL org.opencontainers.image.revision ${REVISION:-unset}
LABEL org.opencontainers.image.base.name scratch

FROM --platform=$BUILDPLATFORM golang:1.23.5-alpine AS checker-build
ARG TARGETOS
ARG TARGETARCH
COPY ./ .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o status-checker main.go

FROM status AS checker
LABEL org.opencontainers.image.ref.name "sbnarra/status-checker"
LABEL org.opencontainers.image.title "status-checker"
LABEL org.opencontainers.image.description "status checker"
COPY ./config /config
COPY --from=checker-build /go/status-checker /status-checker
ENTRYPOINT ["/status-checker"]

EXPOSE 8000
VOLUME /config
VOLUME /history
