FROM golang:1.26.3-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /app/server ./cmd/api

FROM alpine:3.22
WORKDIR /app
COPY --from=build /app/server /app/server
COPY --from=build /app/api/v1 /app/api/v1
COPY --from=build /app/web/swagger /app/web/swagger
COPY --from=build /app/database/migrations /app/database/migrations
EXPOSE 8080
ENTRYPOINT ["/app/server"]
