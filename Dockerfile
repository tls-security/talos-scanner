# --- Estágio 1: Builder (Compila o Go) ---
FROM --platform=linux/amd64 golang:1.23-alpine AS builder

WORKDIR /app

# Instala git para baixar dependências
RUN apk add --no-cache git

# Copia arquivos de dependência
COPY go.mod go.sum ./
RUN go mod download

# Copia o código fonte
COPY . .

# Compila o binário estático
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/api/main.go

# --- Estágio 2: Runner (Imagem Final com Chrome) ---
FROM --platform=linux/amd64 alpine:latest

# Instala Chromium e dependências gráficas
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    dumb-init

# Variáveis para o Chrome
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/lib/chromium/

WORKDIR /app

# Copia o binário e a pasta public
COPY --from=builder /app/server .
COPY --from=builder /app/public ./public

# Configura porta
ENV PORT=8080
EXPOSE 8080

# dumb-init gerencia os processos do Chrome corretamente
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["./server"]