---
if: docker
---
# Build Step
FROM golang:1.27.3-alpine@sha256:dbaa054703263d06c8b49576667418473b6231f8c0fa69dee103207fdb8e3654 as builder

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
FROM gcr.io/distroless/static@sha256:448810357cc8e67b7208988abdba0f5ec3d3aebee0a5b13eb624979792cffc0f
COPY --from=builder /app/{{ .name }} /go/bin/{{ .name }}
ENTRYPOINT ["/go/bin/{{ .name }}"]
