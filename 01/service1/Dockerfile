FROM golang:1.18 as builder
WORKDIR /workspace
COPY . /workspace
RUN CGO_ENABLED=0 GOOS=linux go build -o sample main.go && chmod +x ./sample

FROM alpine:3.15
WORKDIR /
RUN apk --no-cache add ca-certificates
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=builder /workspace/ ./

EXPOSE 8080
CMD ["./sample"]
