FROM golang:1.20-bookworm

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . /app

RUN go build -o /journalink

EXPOSE 8080

CMD [ "/journalink" ]

