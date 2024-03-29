FROM golang:1.20-alpine AS build
LABEL authors="cruii"
WORKDIR /app

ADD . .
RUN go mod download
RUN go build -o /app/bilibili cmd/main.go

FROM alpine:3.14
WORKDIR /app
RUN apk update && apk add tzdata
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" > /etc/timezone
COPY --from=build /app/bilibili /app/
ENTRYPOINT ["./bilibili"]