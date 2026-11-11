# syntax=docker/dockerfile:1.20
FROM golang:1.27.4-bookworm@sha256:4fa7c968ba93f1635670522dbb01a2c78150443e1385746924566079f36570ea AS builder

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
