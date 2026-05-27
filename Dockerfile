FROM golang:1.24.4-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /go-adega-service .

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /go-adega-service /app/go-adega-service
COPY db/migration /app/db/migration

EXPOSE 8085

CMD ["/app/go-adega-service"]
