FROM golang:1.26.2

WORKDIR /test-task

ENV GOPROXY=https://proxy.golang.org,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app main.go

CMD ["./app"]