FROM golang:1.21

WORKDIR /app

COPY . .

RUN go mod download

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o lunchmoney-server ./cmd/

CMD ["/app/lunchmoney-server"]