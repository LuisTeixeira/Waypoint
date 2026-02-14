#-----Build-----
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o waypoint-api ./cmd/server/main.go

#-----Run-----
FROM alpine:3.21

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

WORKDIR /app

COPY --from=builder /app/waypoint-api /waypoint-api

EXPOSE 8080

ENTRYPOINT [ "/waypoint-api" ]