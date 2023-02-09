FROM golang:1.18-alpine AS builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
# COPY go.mod go.sum ./
# RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/emitter .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /usr/local/bin/emitter /usr/local/bin/emitter
COPY --from=builder /usr/src/app/text*.txt /root/

EXPOSE 9988/tcp

CMD ["emitter"]