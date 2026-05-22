# 1. Use the absolute latest Go version to match your local machine
FROM golang:alpine

WORKDIR /app

# 2. Use the updated 'air-verse' repository!
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air"]