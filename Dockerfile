# Multi-stage build for a small final image
FROM golang:1.21-alpine AS builder

WORKDIR /src
COPY go.mod .
COPY . .

# Build static binary
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app .

FROM gcr.io/distroless/static:nonroot

ENV PORT=8080
WORKDIR /
COPY --from=builder /out/app /app

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app"]

