FROM golang:1.22.3-alpine as builder

WORKDIR /app

COPY ../ ./

RUN GOFLAGS="-buildvcs=false" go generate ./...
RUN GOFLAGS="-buildvcs=false" go build -o coma .

FROM alpine:latest

RUN apk update && \
    apk --update add bash

WORKDIR /app 

EXPOSE 5899

COPY --from=builder /app/coma /app

CMD /app/coma