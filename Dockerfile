FROM --platform=$BUILDPLATFORM golang:1.23.5-alpine AS go
ARG TARGETOS
ARG TARGETARCH
ARG MODULE

WORKDIR /
COPY ./ .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o status-${MODULE:-checker} cmd/${MODULE:-checker}/main.go

FROM scratch
ARG CREATED
ARG REVISION
ARG MODULE

# https://github.com/opencontainers/image-spec/blob/main/annotations.md
LABEL org.opencontainers.image.ref.name "sbnarra/status-checker"
LABEL org.opencontainers.image.title "status-${MODULE}"
LABEL org.opencontainers.image.description "status ${MODULE}"
LABEL org.opencontainers.image.source "https://github.com/sbnarra/status-checker"
LABEL org.opencontainers.image.documentation "https://github.com/sbnarra/status-checker"
LABEL org.opencontainers.image.created ${CREATED:-unset}
LABEL org.opencontainers.image.version ${REVISION:-unset}
LABEL org.opencontainers.image.revision ${REVISION:-unset}
LABEL org.opencontainers.image.base.name scratch

COPY --from=go /status-${MODULE:-checker} /status-${MODULE:-checker}
ENTRYPOINT ["/status-${MODULE:-checker}"]