FROM golang:1.22.4

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify


COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/test

CMD [ "app" ]
