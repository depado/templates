---
if: docker
---
# Build Step
FROM golang:1.26.5-alpine@sha256:0178a641fbb4858c5f1b48e34bdaabe0350a330a1b1149aabd498d0699ff5fb2 AS builder

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
FROM gcr.io/distroless/static@sha256:d5f030ca7c5793784e9ea4178a116da360250411d13921a5af27c6cb5a5949bf
COPY --from=builder /app/{{ .name }} /app/{{ .name }}
WORKDIR /app
EXPOSE 8090
ENTRYPOINT ["/app/{{ .name }}", "serve", "--http=0.0.0.0:8090"]
