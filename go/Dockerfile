---
if: docker
---
# Build Step
FROM golang:1.24.5-alpine@sha256:9c4b616be9d26e4762219223331bab5db98649e4be1f6badeac3f7c00a340e3f as builder

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
FROM gcr.io/distroless/static@sha256:b7b9a6953e7bed6baaf37329331051d7bdc1b99c885f6dbeb72d75b1baad54f9
COPY --from=builder /app/{{ .name }} /go/bin/{{ .name }}
ENTRYPOINT ["/go/bin/{{ .name }}"]
