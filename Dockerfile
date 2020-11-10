FROM golang:1.14 as builder
RUN mkdir -p /app
WORKDIR /app
COPY . .
RUN go get -d ./cmd/.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o previewer ./cmd/.

FROM alpine:latest
WORKDIR /
COPY --from=builder /app/previewer .
RUN mkdir -p /assets/cache
CMD ["/previewer"]
EXPOSE 8080