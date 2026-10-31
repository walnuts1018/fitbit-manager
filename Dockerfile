# syntax=docker/dockerfile:1.20
FROM golang:1.27.3-bookworm@sha256:d9b6a1beb9474adc2a1ce842573639b44b2a0dd49b6669f51d3821baf9ec37c9 AS builder

ENV ROOT=/build
RUN mkdir ${ROOT}
WORKDIR ${ROOT}

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=linux go build -o fitbit-manager $ROOT/cmd/server && chmod +x ./fitbit-manager

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=linux go build -o fitbit-manager-job $ROOT/cmd/job && chmod +x ./fitbit-manager-job

FROM gcr.io/distroless/cc-debian13:nonroot
WORKDIR /app

COPY  ./templates/ /app/templates/
COPY  ./assets/ /app/assets/

COPY --from=builder /build/fitbit-manager ./
COPY --from=builder /build/fitbit-manager-job ./

CMD ["./fitbit-manager"]
LABEL org.opencontainers.image.source "https://github.com/walnuts1018/fitbit-manager"
