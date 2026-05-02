FROM golang:1.26

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 3000

CMD ["go", "run", "./cmd/main.go"]
