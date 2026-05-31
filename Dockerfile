FROM node:22-alpine AS frontend-builder

WORKDIR /app/frontend
ARG VITE_API_BASE_URL=/api
ARG VITE_MAIN_APP_BASE_URL=
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}
ENV VITE_MAIN_APP_BASE_URL=${VITE_MAIN_APP_BASE_URL}
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile
COPY frontend/ ./
RUN pnpm build

FROM golang:1.24-alpine AS backend-builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY backend/ ./backend
COPY --from=frontend-builder /app/backend/internal/web/dist ./backend/internal/web/dist
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/sub2api-distributor ./backend/cmd/server

FROM alpine:3.22

RUN addgroup -S distributor && adduser -S distributor -G distributor \
  && apk add --no-cache ca-certificates tzdata wget postgresql17-client

WORKDIR /app
COPY --from=backend-builder /app/sub2api-distributor /app/sub2api-distributor
COPY --from=frontend-builder /app/backend/internal/web/dist /app/web
COPY backend/migrations /app/migrations
COPY deploy/docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod +x /app/docker-entrypoint.sh

ENV APP_ENV=production
ENV SERVER_PORT=8091
ENV STATIC_DIR=/app/web
ENV RUN_MIGRATIONS_ON_STARTUP=true

EXPOSE 8091

HEALTHCHECK --interval=30s --timeout=10s --retries=3 --start-period=20s \
  CMD wget -q -T 5 -O /dev/null http://localhost:${SERVER_PORT}/health || exit 1

USER distributor

ENTRYPOINT ["/app/docker-entrypoint.sh"]
