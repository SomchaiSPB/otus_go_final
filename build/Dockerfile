# syntax=docker/dockerfile:1

FROM golang:1.17-alpine AS build_base

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o ./previewer cmd/main.go

FROM alpine:3.9.6

COPY --from=build_base /app/previewer /app/previewer

EXPOSE 4000

CMD ["./app/previewer"]