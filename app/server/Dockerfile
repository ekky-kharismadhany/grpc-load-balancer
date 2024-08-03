FROM golang:1.22-alpine3.20 AS builder

WORKDIR /app

COPY . .

RUN go clean && CGO_ENABLED=0 go build -o main ./main.go

FROM alpine:3.20.2

WORKDIR /app

COPY --from=builder /app/main ./

ENTRYPOINT [ "./main" ]
