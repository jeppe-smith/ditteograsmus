FROM golang:1.22 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server .

FROM scratch
COPY --from=build /server /server
EXPOSE 1337
CMD ["/server"]
