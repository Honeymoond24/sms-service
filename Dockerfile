FROM golang:1.22.5-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/app/


FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/.env .
COPY --from=build /app/db ./db
COPY --from=build /app/internal/infrastructure/database/migrations ./internal/infrastructure/database/migrations
COPY --from=build /app/main .

CMD ["./main"]
