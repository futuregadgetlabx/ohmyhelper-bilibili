FROM golang:1.20-alpine AS build
LABEL authors="cruii"
WORKDIR /app

ADD . .
RUN go mod download
RUN go build -o /app/bilibili cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/bilibili /app/
ENTRYPOINT ["./bilibili"]