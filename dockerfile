# production environment
FROM golang:1.12
# setup go
ENV GOBIN $GOPATH/bin
ENV PATH $GOBIN:/usr/local/go/bin:$PATH
ENV GO11MODULE "auto"
WORKDIR /go/src/dgoldstein1/edge

# build server binary
COPY . /go/src/dgoldstein1/edge
RUN go get 
RUN go build -o edge
RUN ls edge

ENV PORT "8443"
CMD ./edge