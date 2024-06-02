FROM golang:1.22.3-alpine as builder

WORKDIR /app

COPY ../ ./

RUN GOFLAGS="-buildvcs=false" go generate ./...
RUN GOFLAGS="-buildvcs=false" go build -o coma .

FROM alpine:latest

RUN apk update && \
    apk --update add bash nano
    
RUN adduser -D -g '' comauser

WORKDIR /usr/local/opt/coma
WORKDIR /var/lib/coma
COPY --from=builder /app/coma /usr/local/bin/coma

RUN chown comauser:comauser /var/lib/coma \
    && chmod 750 /var/lib/coma
RUN chown comauser:comauser /usr/local/opt/coma \
    && chmod 750 /usr/local/opt/coma

USER comauser

CMD /usr/local/bin/coma