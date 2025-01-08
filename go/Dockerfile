---
if: docker
---
# Build Step
FROM golang:1.23.4-alpine@sha256:94b4686346804a8cfc3d2c9069bbdf4bfb21ac8bdc6d99e23f529e295db1cd81 as builder

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
FROM gcr.io/distroless/static@sha256:5c7e2b465ac6a2a4e5f4f7f722ce43b147dabe87cb21ac6c4007ae5178a1fa58
COPY --from=builder /app/{{ .name }} /go/bin/{{ .name }}
ENTRYPOINT ["/go/bin/{{ .name }}"]
