FROM golang:latest

WORKDIR /app/demo
COPY . .

RUN go build eshop_cart

EXPOSE 8888
ENTRYPOINT ["./eshop_cart"]