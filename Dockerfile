FROM golang:1.25.3-alpine@sha256:91eda9776261207ea25fd06b5b7fed8d397dd2c0a283e77f2ab6e91bfa71079d AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app

FROM alpine:3@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11 AS runtime

WORKDIR /app
COPY --from=build /app/app .

RUN apk add tzdata
ENV TZ=Asia/Tokyo

RUN mkdir -p /app/data

CMD ["./app"]
