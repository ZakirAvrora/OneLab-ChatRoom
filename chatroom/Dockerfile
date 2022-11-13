FROM golang:1.18 as builder
COPY . /
WORKDIR /
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o golang-app .

FROM alpine:3.16
COPY --from=builder /golang-app .
COPY --from=builder /public /public
EXPOSE 8080:8080
ENTRYPOINT ["./golang-app", "-redis-address", "redis-db"]