FROM golang:1.14-alpine AS builder

WORKDIR /go/src/app
COPY . .

RUN go get -v -t -d ./...
RUN go build -v ./cmd/...

FROM alpine

WORKDIR /

COPY --from=builder /go/src/app/unsocket .

ENTRYPOINT ["/unsocket"]
