# syntax=docker/dockerfile:1

FROM golang:1.22

WORKDIR /app

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /server

RUN mkdir -p ./uploads

VOLUME ["/app/uploads"]

EXPOSE 1337

CMD ["/server"]
