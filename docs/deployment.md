# Guia de Implantação (Deployment Guide)

Este guia descreve as configurações e comandos necessários para realizar o build e o deployment desta API Go nos provedores de hospedagem **Render**, **Vercel** e **AWS**.

---

## 1. Requisitos Globais de Ambiente
Independentemente da plataforma, a API necessita das seguintes configurações ambientais:

- **Variável de Ambiente**:
  - `DATABASE_URL`: A string de conexão do PostgreSQL (Neon). Exemplo: `postgresql://usuario:senha@host/neondb?sslmode=require`
- **Porta**:
  - A API escuta na porta `:8080` por padrão. Certifique-se de expor ou mapear essa porta se o host não injetar a variável de ambiente `PORT` (ou você pode opcionalmente ajustar no código para ler `os.Getenv("PORT")`).

---

## 2. Implantação no Render (Recomendado para APIs Go)

O **Render** oferece suporte nativo a Go através de Web Services.

### Configuração:
1. Acesse o painel do Render e crie um novo **Web Service**.
2. Conecte o seu repositório Git.
3. Configure os campos de Build e Start:
   - **Runtime**: `Go`
   - **Build Command**:
     ```bash
     go build -o bin/api ./cmd/api
     ```
   - **Start Command**:
     ```bash
     ./bin/api
     ```
4. Em **Environment Variables**, adicione:
   - `DATABASE_URL` = *(Sua connection string do Neon)*
   - `PORT` = `8080` (Opcional, se precisar mapear no proxy do Render)

---

## 3. Implantação na AWS (Dockerizado)

Para a AWS, a melhor abordagem e mais escalável (seja via **AWS App Runner**, **AWS ECS/Fargate** ou **AWS Elastic Beanstalk**) é usar um container Docker.

### Dockerfile de Produção Recomendado
Crie um arquivo chamado `Dockerfile` na raiz do projeto:

```dockerfile
# Estágio de Build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o api ./cmd/api

# Estágio de Execução
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/api .
EXPOSE 8080
ENV PORT=8080
CMD ["./api"]
```

### Serviços AWS:

#### A. AWS App Runner (Mais simples e moderno)
1. Conecte o repositório no App Runner.
2. Escolha **Branch deployment** e configure o método de build para **Dockerfile** (ou utilize o arquivo `Dockerfile` criado acima).
3. Adicione `DATABASE_URL` nas variáveis de ambiente.
4. Defina a porta de escuta para `8080`.

#### B. AWS Elastic Beanstalk
1. Crie uma aplicação Elastic Beanstalk selecionando a plataforma **Docker**.
2. Suba o código contendo o `Dockerfile` na raiz.
3. Nas propriedades de software da instância, configure a variável de ambiente `DATABASE_URL`.

---

## 4. Implantação na Vercel (Serverless)

A **Vercel** é voltada para arquiteturas serverless/JAMstack. Para rodar Go nativamente na Vercel, usamos o runtime oficial/comunidade para Serverless Functions.

### Configuração:
1. Crie um diretório chamado `api/` na raiz do projeto (a Vercel mapeia arquivos nessa pasta automaticamente como funções serverless).
2. Adicione um arquivo `vercel.json` na raiz do projeto para rotear as requisições para os arquivos Go serverless correspondentes:

#### Exemplo de `vercel.json`:
```json
{
  "version": 2,
  "builds": [
    {
      "src": "api/**/*.go",
      "use": "@vercel/go"
    }
  ],
  "routes": [
    {
      "src": "/(.*)",
      "dest": "api/main.go"
    }
  ]
}
```

#### Exemplo de Handler Serverless em `api/main.go`:
Para rodar como Serverless Function, o arquivo de entrada Go deve expor a assinatura do `net/http` padrão da Vercel:

```go
package handler

import (
	"net/http"
	// importe o roteador do seu projeto se necessário para despachar as rotas
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Ponto de entrada que delega para o seu mux / roteador
}
```
3. Defina `DATABASE_URL` no painel da Vercel nas configurações de **Environment Variables** do projeto.
