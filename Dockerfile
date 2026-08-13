# ---- Etapa 1: build ----
FROM golang:1.25-alpine AS build

WORKDIR /app

# Copiamos solo los archivos de dependencias primero para aprovechar la cache de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto del código
COPY . .

# Compilamos el binario:
# - CGO_ENABLED=0 -> binario estático, sin dependencias de libc (necesario para distroless)
# - GOOS=linux    -> por si se construye desde Mac/Windows
# - -ldflags="-s -w" -> quita símbolos de debug, reduce tamaño del binario
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/api

# ---- Etapa 2: imagen final ----
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copiamos únicamente el binario compilado (nada de código fuente, nada de Go toolchain)
COPY --from=build /app/server /server

# Cloud Run inyecta la variable PORT automáticamente (normalmente 8080)
EXPOSE 8080

ENTRYPOINT ["/server"]
