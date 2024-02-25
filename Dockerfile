FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o tp-proxy ./cmd/proxy/main.go

CMD ["./tp-proxy"]