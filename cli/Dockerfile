FROM golang:1.18-alpine AS builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
# COPY go.mod go.sum ./
# RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/word-cli .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /usr/local/bin/word-cli /usr/local/bin/word-cli

CMD ["exec /bin/sh -c 'trap : TERM INT; sleep infinity & wait'"]