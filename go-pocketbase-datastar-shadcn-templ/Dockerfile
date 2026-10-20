---
if: docker
---
# Build Step
FROM golang:1.27.2-alpine@sha256:f92b6ef800e499660581efdabdf25d9d817a9d124eaf900924f0504e7e27e12d AS builder

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
FROM gcr.io/distroless/static@sha256:9e0130bec991f471279b70f1630ada11e46dad439ef2949544f96001225329fb
COPY --from=builder /app/{{.name}} /app/{{.name}}
WORKDIR /app
EXPOSE 8090
ENTRYPOINT ["/app/{{.name}}", "serve", "--http=0.0.0.0:8090"]
