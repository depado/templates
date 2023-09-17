---
if: docker
---
# Build Step
FROM golang:1.21-alpine AS builder

# Dependencies
RUN apk update && apk add --no-cache make git

# Source
WORKDIR $GOPATH/src/{{ .gitserver }}/{{ .owner }}/{{ .name }}
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
COPY . .

# Build
RUN make tmp

# Final Step
FROM gcr.io/distroless/static
COPY --from=builder /tmp/{{ .name }} /go/bin/{{ .name }}
ENTRYPOINT ["/go/bin/{{ .name }}"]
