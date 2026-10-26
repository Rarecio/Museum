# Build stage
FROM golang:1.25-alpine AS builder

LABEL stage=gobuilder

WORKDIR /app

RUN apk update --no-cache && apk add --no-cache tzdata

COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /museum .


# Run stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /museum .

EXPOSE 8080

ENTRYPOINT ["/app/museum"]
