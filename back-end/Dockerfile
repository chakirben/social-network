FROM golang:1.23.0-alpine as builder

ENV CGO_ENABLED=1

RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o main ./server
RUN ls -l /app
FROM alpine:3.18 

RUN apk add --no-cache sqlite-libs
WORKDIR /app
COPY --from=builder /app/main .

COPY --from=builder /app/dataBase/migrations /app/dataBase/migrations



EXPOSE 8080 

CMD ["./main"]