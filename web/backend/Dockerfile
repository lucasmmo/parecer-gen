FROM golang:latest AS build

WORKDIR /app

COPY go.* .

RUN go mod tidy

COPY . .

RUN go build -o /web ./cmd/main.go

FROM ubuntu:jammy AS run

RUN apt-get update && apt-get install -y wkhtmltopdf
ENV LANG en_US.utf8

COPY --from=build /web /web

WORKDIR /app
EXPOSE 8080
CMD ["/web"]