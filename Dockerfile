FROM golang:1.21.4-alpine3.18 as builder
WORKDIR /app
ENV GOOS=linux
ENV GOARCH=amd64
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o sms2ntfy

FROM alpine:3.18.4
WORKDIR /app
COPY --from=builder /app/sms2ntfy /app/
RUN chmod +x /app/sms2ntfy
EXPOSE 8080
ENTRYPOINT [ "/app/sms2ntfy" ]