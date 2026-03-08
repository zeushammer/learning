FROM golang:1.26.1
WORKDIR /app
COPY . .
RUN go build -o server main.go
EXPOSE 1337
CMD ["./server"]
