FROM golang:1.17-alpine as build

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main .

FROM scratch as final

COPY --from=build /app/main /usr/bin/
ENV GIN_MODE release
EXPOSE 80
CMD ["main"]