# production environment
FROM golang:1.12
# setup go
ENV GOBIN $GOPATH/bin
ENV PATH $GOBIN:/usr/local/go/bin:$PATH
ENV GO11MODULE "auto"
WORKDIR /go/src/dgoldstein1/reverseProxy

# build server binary
COPY . /go/src/dgoldstein1/reverseProxy
RUN go get 
RUN go build -o reverseProxy
RUN ls reverseProxy

ENV PORT "8443"
CMD ./reverseProxy