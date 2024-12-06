FROM golang:1.23-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
RUN chmod +x ./wait-for.sh ./start.sh
COPY db/migration ./db/migration

EXPOSE 8080 9090
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]