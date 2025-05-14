FROM golang:1.19

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o main .

EXPOSE 5050

CMD ["./main"]