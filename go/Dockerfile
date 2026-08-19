---
if: docker
---
# Build Step
FROM golang:1.27.0-alpine@sha256:7d5cbf6833f7331dafd25a2e8b9673477f559759ff8ed4ca8efabe6795ad08db as builder

# Dependencies
RUN apk add --no-cache make git

# Source
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
COPY . .

# Build
RUN make

# Final Step
FROM gcr.io/distroless/static@sha256:9197324ba51d9cd071af8505989365c006adf9d6d2067eada25aef00abbb5278
COPY --from=builder /app/{{ .name }} /go/bin/{{ .name }}
ENTRYPOINT ["/go/bin/{{ .name }}"]
