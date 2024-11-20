FROM golang:1.23.3-alpine

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . /app

RUN go build -o main ./main.go

EXPOSE 8081

CMD ["./main"]
