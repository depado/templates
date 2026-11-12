---
if: docker
---
# Build Step
FROM golang:1.27.4-alpine@sha256:0a5f7fe796599ff487605987432e7619ae45beeb2fba4c2ce0b907879932aebf AS builder

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
FROM gcr.io/distroless/static@sha256:1a64f46364620465ef9e9a8f5413757dfdf22695e050bbfd7303eeaf17a59383
COPY --from=builder /app/{{.name}} /app/{{.name}}
WORKDIR /app
EXPOSE 8090
ENTRYPOINT ["/app/{{.name}}", "serve", "--http=0.0.0.0:8090"]
