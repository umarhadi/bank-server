FROM golang:1.19.4-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
RUN chmod +x ./wait-for.sh ./start.sh
COPY db/migration ./db/migration

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]