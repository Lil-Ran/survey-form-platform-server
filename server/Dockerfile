FROM golang:latest AS builder

WORKDIR /app
COPY . .

ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .


FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .

CMD ["./main"]
