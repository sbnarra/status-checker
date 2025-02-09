FROM --platform=$BUILDPLATFORM golang:1.23.5-alpine AS builder
ARG TARGETOS
ARG TARGETARCH
COPY ./ .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o status-checker main.go

FROM scratch
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
LABEL org.opencontainers.image.ref.name "sbnarra/status-checker"
LABEL org.opencontainers.image.title "status-checker"
LABEL org.opencontainers.image.description "status checker"

COPY --from=builder /go/status-checker /status-checker
COPY ./ui /ui
COPY ./config /config

VOLUME /config
VOLUME /history

EXPOSE 8000
ENTRYPOINT ["/status-checker"]