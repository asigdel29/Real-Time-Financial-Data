FROM golang:latest AS build

WORKDIR /app
COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

FROM alpine:latest

RUN apk update && \
    apk upgrade && \
    apk add ca-certificates

WORKDIR /

EXPOSE 8080

COPY --from=build /app/main ./

RUN env && pwd && find .

CMD ["./main"]
