FROM golang:1.26.3-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/gateway-service ./gateway-service/cmd/gateway

FROM alpine:3.22
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /out/gateway-service /app/gateway-service
EXPOSE 8081
ENTRYPOINT ["/app/gateway-service"]
