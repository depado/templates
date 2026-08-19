---
if: docker
---
# Build Step
FROM golang:1.27.0-alpine@sha256:7d5cbf6833f7331dafd25a2e8b9673477f559759ff8ed4ca8efabe6795ad08db AS builder

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
FROM gcr.io/distroless/static@sha256:9197324ba51d9cd071af8505989365c006adf9d6d2067eada25aef00abbb5278
COPY --from=builder /app/{{ .name }} /app/{{ .name }}
WORKDIR /app
EXPOSE 8090
ENTRYPOINT ["/app/{{ .name }}", "serve", "--http=0.0.0.0:8090"]
