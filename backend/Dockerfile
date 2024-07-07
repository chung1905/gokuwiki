FROM golang:alpine as builder

WORKDIR /app
COPY . .
RUN go build .

FROM golang:alpine
WORKDIR /app
COPY --from=builder /app/gokuwiki ./gokuwiki
CMD ["./gokuwiki"]