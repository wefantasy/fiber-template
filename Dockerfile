FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o app .

FROM alpine:latest
COPY --from=builder /app/app /app/app
ENV PORT 8888
ENTRYPOINT ["/app/app"]
EXPOSE 8888