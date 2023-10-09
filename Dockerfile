FROM golang:1.21.2-alpine3.18 AS build

COPY ./ /app

WORKDIR /app

RUN go get ./...
RUN GOOS=linux GOARCH=amd64 go build -o main -buildvcs=false ./cmd/milky-mailer/main.go

FROM scratch

COPY --from=build /app/main /main
COPY --from=build /etc/ssl/certs/ /etc/ssl/certs/

CMD ["./main"]
