FROM golang:1.24-bookworm

WORKDIR /app

# копируем модули и vendor внутрь
COPY go.mod go.sum ./
COPY vendor ./vendor

# говорим го не ходить в интернет, а юзать vendor
ENV GO111MODULE=on
ENV GOPROXY=off
ENV GONOSUMDB=*
ENV GOMODCACHE=/app/vendor

# копируем весь остальной код
COPY . .

# билдим
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o server ./cmd/server

RUN mkdir -p /app/storage-data

EXPOSE 8080
CMD ["/app/server"]
