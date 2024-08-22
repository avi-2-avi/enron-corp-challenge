FROM golang:1.22-alpine AS build-stage

# Building the app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

ARG ENV_FILE=.env
COPY ${ENV_FILE} .env

RUN go build -o main ./cmd/backend/main.go

# Running the app
FROM alpine:latest AS production-stage

WORKDIR /root/

COPY --from=build-stage /app/main .
COPY --from=build-stage /app/.env .env

EXPOSE 8080

CMD ["./main"]