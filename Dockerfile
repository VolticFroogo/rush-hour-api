FROM golang:1.15

WORKDIR /go/src/github.com/VolticFroogo/rush-hour-api
COPY . .
RUN go build -o main .

CMD ["./main"]
