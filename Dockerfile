# Choose whatever you want, version >= 1.16
FROM golang:1.21-alpine

WORKDIR /app
RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]
