FROM golang:1.20.4 as builder
WORKDIR /workspace
COPY . /workspace
RUN CGO_ENABLED=0 GOOS=linux go build -o todoapi main.go && chmod +x ./todoapi

FROM alpine:3.18
WORKDIR /app
RUN apk --no-cache add ca-certificates
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=builder /workspace/ ./

EXPOSE 8080
CMD ["./todoapi"]