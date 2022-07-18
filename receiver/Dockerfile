FROM golang:1.18

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
# COPY go.mod go.sum ./
# RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/receiver .

EXPOSE 9988/tcp
EXPOSE 8899/tcp

CMD ["receiver"]