# Etapa 1: compilar el binario con la imagen oficial de Go
FROM golang:1.27-alpine AS build
WORKDIR /app

# Primero solo las dependencias, para aprovechar la caché de Docker
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o api-go .

# Etapa 2: imagen final liviana, solo con el binario
FROM alpine:3.22
WORKDIR /app
COPY --from=build /app/api-go .

# No correr como root
RUN adduser -D app
USER app

EXPOSE 8080
CMD ["./api-go"]