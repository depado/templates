---
if: docker
---
# Build Step
FROM golang:1.25.1-alpine@sha256:26264873df6ac01d9a0bbff432c83d11572b323efb3bac97c620fff7af4e9fbe as builder

# Dependencies
RUN apk update && apk add --no-cache make git

# Source
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
COPY . .

# Build
RUN make

# Final Step
FROM gcr.io/distroless/static@sha256:f2ff10a709b0fd153997059b698ada702e4870745b6077eff03a5f4850ca91b6
COPY --from=builder /app/{{ .name }} /go/bin/{{ .name }}
ENTRYPOINT ["/go/bin/{{ .name }}"]
