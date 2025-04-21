# Build stage for Node.js/Tailwind
FROM node:20-alpine AS tailwind-builder
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install
COPY tailwind.config.js ./
COPY templates/ ./templates/
RUN npx tailwindcss -i ./input.css -o ./css/output.css --minify

# Build stage for Go
FROM golang:1.22.3-alpine AS go-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY --from=tailwind-builder /app/css/output.css ./css/output.css
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping ./cmd/main

# Final lightweight image
FROM alpine:latest
WORKDIR /
COPY --from=go-builder /docker-gs-ping /docker-gs-ping
COPY --from=tailwind-builder /app/css/output.css ./css/output.css
COPY css/ ./css/
COPY templates/ ./templates/

EXPOSE 8080
CMD ["/docker-gs-ping"]
