FROM golang:1.24.5

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o payment_system ./cmd

CMD ["/app/payment_system"]