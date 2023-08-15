FROM golang:1.20-alpine
LABEL authors="cruii"

ENTRYPOINT ["top", "-b"]