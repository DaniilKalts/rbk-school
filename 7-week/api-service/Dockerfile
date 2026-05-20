FROM golang:1.26.3-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /app/api-service ./cmd/api

FROM alpine:3.22
WORKDIR /app
COPY --from=build /app/api-service /app/api-service
COPY --from=build /app/api-service/api/v1 /app/api-service/api/v1
COPY --from=build /app/api-service/web/swagger /app/api-service/web/swagger
COPY --from=build /app/api-service/database/migrations /app/api-service/database/migrations
EXPOSE 8080
ENTRYPOINT ["/app/api-service"]
