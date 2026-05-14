FROM golang:1.26

ENV TZ=Africa/Addis_Ababa
RUN apt-get update && apt-get install -y tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/air-verse/air@latest

COPY . .

EXPOSE 3000

CMD ["air"]
