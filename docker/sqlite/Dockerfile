FROM golang:1.24-alpine AS build

# Install build dependencies required for CGO/SQLite
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY ../.. .

# Explicitly enable CGO and build (Gorm Sqlite drivers are based on CGO)
ENV CGO_ENABLED=1
RUN go build -o main cmd/api/main.go

FROM alpine:3.21.3 AS prod

# SQLite runtime dependencies
RUN apk add --no-cache sqlite-libs

WORKDIR /app
COPY --from=build /app/main /app/main

ARG PORT
EXPOSE ${PORT}

CMD ["./main"]