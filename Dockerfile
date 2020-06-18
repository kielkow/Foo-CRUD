FROM golang:1.14

WORKDIR /go/src/

COPY . .

RUN GOOS=linux go build -ldflags="-s -w"

EXPOSE 3000

ENTRYPOINT ["./main"]