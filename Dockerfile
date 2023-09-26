# Use a imagem oficial do Golang como base
FROM golang:latest

# Defina a variável de ambiente PORT com um valor padrão


# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copie o código fonte para o diretório de trabalho
COPY . .

# Compile a aplicação
RUN go build ./cmd/main.go

# Expor a porta configurada
EXPOSE ${PORT}

# Comando para iniciar a aplicação
CMD ["./main"]
