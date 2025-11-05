# ---------- build stage ----------
FROM golang:1.24-bookworm AS builder

WORKDIR /app

# модули
COPY backend/go.mod backend/go.sum ./
COPY backend/vendor ./vendor

ENV GO111MODULE=on
ENV GOPROXY=off
ENV GOMODCACHE=/app/vendor

# весь код
COPY backend/ ./

# собрали бинарь
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o /app/server ./cmd/server

# тут же создадим пустые директории, чтобы потом просто скопировать
RUN mkdir -p /app/storage-data /app/data /app/secrets


# ---------- final stage ----------
FROM gcr.io/distroless/static-debian12:latest

WORKDIR /app

# копируем бинарь и заранее созданные папки из билдера
COPY --from=builder /app/server /app/server
COPY --from=builder /app/storage-data /app/storage-data
COPY --from=builder /app/secrets /app/secrets

EXPOSE 8080
CMD ["/app/server"]
