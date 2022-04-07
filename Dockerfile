FROM golang:1.17.8-alpine3.15

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o /open-banking

ENV GIN_MODE=release

CMD [ "/open-banking" ]
