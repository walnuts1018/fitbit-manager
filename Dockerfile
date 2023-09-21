FROM golang:1.21 as builder
ENV ROOT=/build
RUN mkdir ${ROOT}
WORKDIR ${ROOT}

COPY ./ ./
RUN go get

RUN CGO_ENABLED=0 GOOS=linux go build -o fitbit-manager $ROOT/main.go && chmod +x ./fitbit-manager

FROM alpine:3
WORKDIR /app

COPY --from=builder /build/fitbit-manager ./
COPY --from=builder /build/templates/ /app/templates/
COPY --from=builder /build/assets/ /app/assets/
COPY --from=builder /usr/share/zoneinfo/Asia/Tokyo /usr/share/zoneinfo/Asia/Tokyo
CMD ["./fitbit-manager"]
LABEL org.opencontainers.image.source = "https://github.com/walnuts1018/fitbit-manager"
