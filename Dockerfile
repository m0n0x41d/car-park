FROM golang:1.19.3-alpine3.16 as builder

COPY . /app
WORKDIR /app
RUN go get -d -v .
RUN go build

FROM alpine:3.16.2
WORKDIR /root/
COPY --from=builder /app ./
CMD ["./car-park"]
EXPOSE 8888