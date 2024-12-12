FROM golang:1.23.4 AS builder
ENV ROOT=/build
RUN mkdir ${ROOT}
WORKDIR ${ROOT}

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o fitbit-manager /$ROOT && chmod +x ./fitbit-manager

FROM debian:bookworm-slim
WORKDIR /app

COPY --from=builder /build/fitbit-manager ./
COPY --from=builder /build/templates/ /app/templates/
COPY --from=builder /build/assets/ /app/assets/
COPY --from=builder /usr/share/zoneinfo/Asia/Tokyo /usr/share/zoneinfo/Asia/Tokyo

CMD ["./fitbit-manager"]
LABEL org.opencontainers.image.source = "https://github.com/walnuts1018/fitbit-manager"
