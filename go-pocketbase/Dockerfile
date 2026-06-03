---
if: docker
---
# Build Step
FROM golang:1.26.4-alpine@sha256:162b257ae1156df92d48e2de5abe69f657fc4961f0ea0e5376837a9ffed8160b AS builder

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
FROM gcr.io/distroless/static@sha256:3592aa8171c77482f62bbc4164e6a2d141c6122554ace66e5cc910cadb961ff0
COPY --from=builder /app/{{ .name }} /app/{{ .name }}
WORKDIR /app
EXPOSE 8090
ENTRYPOINT ["/app/{{ .name }}", "serve", "--http=0.0.0.0:8090"]
