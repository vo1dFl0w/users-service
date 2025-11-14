FROM golang:1.24.3-alpine AS builder

WORKDIR /users-service

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

ENV CONFIG_PATH=/users-service/config/config.yaml
ENV CGO_ENABLED=0

RUN go build -ldflags="-s -w" -o users-service ./cmd/users-service/main.go

FROM alpine AS runner

RUN apk add --no-cache ca-certificates

WORKDIR /users-service

COPY --from=builder /users-service/users-service /users-service/users-service
COPY --from=builder /users-service/config/config.yaml /users-service/config/config.yaml

EXPOSE 8080

CMD ["./users-service"]