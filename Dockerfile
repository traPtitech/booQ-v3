FROM golang:1.25.3-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app

FROM alpine:3 AS runtime

WORKDIR /app
COPY --from=build /app/app .

RUN apk add tzdata
ENV TZ=Asia/Tokyo

RUN mkdir -p /app/data

CMD ["./app"]
