FROM golang:1.23.0-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app

FROM alpine:3 AS runtime

WORKDIR /app
COPY --from=build /app/app .

RUN mkdir -p /app/data

CMD ["./app"]
