FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o message-gateway-service

EXPOSE 8081

CMD ["./message-gateway-service"]
