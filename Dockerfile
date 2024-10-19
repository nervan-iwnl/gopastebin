FROM golang:1.23-alpine

WORKDIR /app

RUN apk add --no-cache curl bash

ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

COPY go.mod go.sum ./

ENV GOTOOLCHAIN=auto

RUN go mod download
RUN go mod tidy

COPY . .

RUN go build -o app .

EXPOSE 8080

CMD ["/wait-for-it.sh", "postgres:5432", "--", "./app"]
