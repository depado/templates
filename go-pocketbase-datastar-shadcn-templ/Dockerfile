---
if: docker
---
# Build Step
FROM golang:1.27.1-alpine@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414 AS builder

# Dependencies
RUN apk add --no-cache make git

# Source
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
COPY . .

# Build
RUN make build

# Final Step
FROM gcr.io/distroless/static@sha256:58133991db06659feaabe0f4e97a35cebf15ef4ea08f8a4c6d2ee5f75e4aa6a0
COPY --from=builder /app/{{.name}} /app/{{.name}}
WORKDIR /app
EXPOSE 8090
ENTRYPOINT ["/app/{{.name}}", "serve", "--http=0.0.0.0:8090"]
