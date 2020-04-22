FROM golang:1.13-alpine

# Install system dependencies
RUN apk add libc-dev=0.7.2-r0
RUN apk add gcc=9.2.0-r4
RUN apk add libpcap-dev=1.9.1-r0

WORKDIR /go/src/inkle
COPY . .

# Get module dependencies
RUN go get .
RUN go build -v .

ENTRYPOINT ["./inkle"]
