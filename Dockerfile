FROM golang:1.18-alpine as gateway
WORKDIR /gateway

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -o gateway ./cmd/gateway/main.go

EXPOSE 8080
CMD ["./gateway"]