# Build stage for Node.js/Tailwind
FROM node:20-alpine AS tailwind-builder
WORKDIR /app
COPY package.json package-lock.json tailwind.config.js ./
COPY input.css .
COPY templates/ ./templates/
RUN npm install
RUN npx tailwindcss -i input.css -o ./static/css/output.css --minify

# Build stage for Go
FROM golang:1.22.3-alpine AS go-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY --from=tailwind-builder /app/static/css/output.css ./static/css/output.css
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /web-proj ./cmd/main

# Final lightweight image
FROM alpine:latest
WORKDIR /
COPY --from=go-builder /web-proj /web-proj
COPY --from=tailwind-builder /app/static/ ./static/
COPY templates/ ./templates/

EXPOSE 8080
CMD ["/web-proj"]
