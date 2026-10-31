# syntax=docker/dockerfile:1.20@sha256:26147acbda4f14c5add9946e2fd2ed543fc402884fd75146bd342a7f6271dc1d
FROM golang:1.25.7-bookworm@sha256:564e366a28ad1d70f460a2b97d1d299a562f08707eb0ecb24b659e5bd6c108e1 AS builder

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

FROM gcr.io/distroless/cc-debian13:nonroot@sha256:3ee00291cf386ec52a552e6a858d555debd74776eb4bfa53fa3df2cd6677515e
WORKDIR /app

COPY  ./templates/ /app/templates/
COPY  ./assets/ /app/assets/

COPY --from=builder /build/fitbit-manager ./
COPY --from=builder /build/fitbit-manager-job ./

CMD ["./fitbit-manager"]
LABEL org.opencontainers.image.source "https://github.com/walnuts1018/fitbit-manager"
