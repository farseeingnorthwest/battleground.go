FROM golang:1.21-alpine AS builder

WORKDIR /opt/src
COPY . .
RUN go build -o battleground

FROM scratch

WORKDIR /opt/bin
COPY --from=builder /opt/src/battleground .
ENTRYPOINT ["/opt/bin/battleground"]

EXPOSE 3000
