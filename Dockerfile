FROM golang:1.22-alpine

WORKDIR /app
COPY . .

RUN go build -o hasher ./main.go

CMD ["./hasher", ":8000"]