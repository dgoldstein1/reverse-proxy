FROM golang:1.22
ENV GO111MODULE=on
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o reverseProxy

ENV PORT="8443"
CMD ["./reverseProxy"]
