FROM golang:1.20-alpine AS build
LABEL authors="cruii"
WORKDIR /app
ENV BILIBILI_USERID="RUNTIME_BILIBILI_USERID" \
    SECRET_KEY="YOUR_SECRET_KEY" \
    AUTHOR_ID=287969457
ADD . .
RUN go mod download
RUN go build -o /app/bilibili cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/bilibili /app/
COPY --from=build /app/conf/config.yaml /app/conf/
ENTRYPOINT ["./bilibili"]