FROM golang:1.18-alpine AS builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
# COPY go.mod go.sum ./
# RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/receiver .


FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /usr/local/bin/receiver /usr/local/bin/receiver

EXPOSE 9988/tcp
EXPOSE 8899/tcp

CMD ["receiver"]