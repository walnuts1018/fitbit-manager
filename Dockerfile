# syntax=docker/dockerfile:1.20
FROM golang:1.25.7-bookworm AS builder

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

FROM gcr.io/distroless/cc-debian13:nonroot@sha256:29ee70762214fe8190cea04ef1c0a84a8365879e047ce8f911ae5a5ab2355fa2
WORKDIR /app

COPY  ./templates/ /app/templates/
COPY  ./assets/ /app/assets/

COPY --from=builder /build/fitbit-manager ./
COPY --from=builder /build/fitbit-manager-job ./

CMD ["./fitbit-manager"]
LABEL org.opencontainers.image.source "https://github.com/walnuts1018/fitbit-manager"
