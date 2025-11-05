# ---------- build ----------
FROM node:20-alpine AS build

WORKDIR /app

# копируем ВСЁ, чтобы не ловить "package.json not found"
COPY . .

# если у тебя pnpm/yarn — поменяешь тут
RUN npm install

# это важно: внутри докер-сетки фронт ходит к сервису "app"
ARG VITE_API_BASE=http://localhost:8080/api/v1

ENV VITE_API_BASE=$VITE_API_BASE

RUN npm run build

# ---------- run ----------
FROM nginx:1.27-alpine

RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
