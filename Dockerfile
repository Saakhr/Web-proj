FROM golang:1.23.4
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
 
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping /app/cmd/main
# RUN openssl genrsa -out /app/key.pem 2048

EXPOSE 8080

CMD ["/docker-gs-ping"]
