FROM golang:1.22.5 as builder

WORKDIR /app

COPY . . 

RUN go mod download

RUN go test -v ./...


RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

FROM scratch as final

EXPOSE 8080

COPY --from=builder /app/main /app/main
COPY --from=builder /app/.env /app/.env

CMD ["/app/main"]