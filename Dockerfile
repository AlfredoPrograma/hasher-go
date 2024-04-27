FROM golang:1.22-alpine

# Install dependencies for CGO
RUN apk add --no-cache --update gcc g++

WORKDIR /app
COPY . .


RUN go get -v
RUN CGO_ENABLED=1 go build -o hasher ./main.go

CMD ["./hasher", ":8000"]