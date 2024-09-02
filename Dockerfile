FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY main.go main.go
COPY internal/ internal/

RUN CGO_ENABLED=0 GOOS=linux go build -o /abc-trading

EXPOSE 8080

CMD ["/abc-trading"]
