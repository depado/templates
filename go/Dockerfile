---
if: docker
---
# Build Step
FROM golang:1.26.6-alpine@sha256:3889b425f035be855a72fb4755265311293b6d414521f0a519d819df32222d83 as builder

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
