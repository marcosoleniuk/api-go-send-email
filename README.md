# 📧 API Go Send Email

[![Go Version](https://img.shields.io/badge/Go-1.22+-blue)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-ready-blue)](https://www.docker.com/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

API desenvolvida em **Golang** para envio de emails com:

- ✅ Corpo do email em **HTML** (tags, links, imagens, etc.)  
- ✅ Suporte a **anexos**  
- ✅ Estrutura modular (config, handlers, middleware, models)  
- ✅ Fácil deploy com **Docker**  
- ✅ Configuração via **variáveis de ambiente (.env)**  

---

## 📂 Estrutura do Projeto

```
api-go-send-email/
│── config/        # Configurações (carregamento de env, smtp, etc.)
│── database/      # (Opcional) Persistência de dados ou logs
│── handlers/      # Handlers HTTP com as rotas
│── middleware/    # Middlewares (ex: autenticação, logs, CORS)
│── models/        # Definições de modelos (request/response)
│── .dockerignore  # Arquivos ignorados no build Docker
│── .env.example   # Exemplo de variáveis de ambiente
│── .gitignore     # Arquivos ignorados no Git
│── Dockerfile     # Containerização da aplicação
│── go.mod         # Gerenciamento de dependências
│── go.sum         # Checksum das dependências
│── main.go        # Ponto de entrada da aplicação
```

---

## ⚙️ Configuração

### Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto, baseado no `.env.example`:

```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=465
SMTP_USER=<seu-email>
SMTP_PASS=<senha-email>

API_USERNAME=<usuário>
API_PASSWORD=<password>
API_KEY=<api-key-gerada>

DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=tb_api_emails

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=sua-senha
```

---

## ▶️ Como Rodar

### 1. Localmente

```bash
# Clone o repositório
git clone https://github.com/marcosoleniuk/api-go-send-email.git
cd api-go-send-email

# Instale dependências
go mod tidy

# Execute a aplicação
go run main.go
```

A API ficará disponível em:

```
http://localhost:8080
```

---

### 2. Com Docker

```bash
# Build da imagem
docker build -t api-go-send-email .

# Rodar o container
docker run -d -p 8080:8080 --env-file .env api-go-send-email
```

---

## 📡 Endpoints Disponíveis

### **1. Enviar Email**
`POST /send-email`

#### Request Body
```json
{
  "to": ["destinatario@dominio.com"],
  "subject": "Teste de Envio",
  "body": "<h1>Olá!</h1><p>Esse é um email de teste com <b>HTML</b>.</p>",
  "attachments": ["uploads/relatorio.pdf"]
}
```

- **to** → lista de destinatários  
- **subject** → assunto do email  
- **body** → corpo em HTML  
- **attachments** → lista de arquivos para anexar (opcional)  

#### Resposta de Sucesso
```json
{
  "message": "Email enviado com sucesso!"
}
```

#### Resposta de Erro
```json
{
  "error": "Falha ao enviar email"
}
```

---

## 🧪 Testes com cURL / HTTPie

### Enviar Email Simples
```bash
curl -X POST http://localhost:8080/send-email   -H "Content-Type: application/json"   -d '{
    "to": ["teste@dominio.com"],
    "subject": "API Go Send Email",
    "body": "<h2>Email enviado com sucesso 🚀</h2>"
  }'
```

### Enviar com Anexo
```bash
curl -X POST http://localhost:8080/send-email   -H "Content-Type: application/json"   -d '{
    "to": ["teste@dominio.com"],
    "subject": "Relatório Anexo",
    "body": "<p>Segue o relatório em anexo.</p>",
    "attachments": ["uploads/relatorio.pdf"]
  }'
```

---

## 🛠 Tecnologias Utilizadas

- [Go](https://go.dev/) – Linguagem principal  
- [net/smtp](https://pkg.go.dev/net/smtp) – Envio de emails  
- [Docker](https://www.docker.com/) – Containerização  
- [godotenv](https://github.com/joho/godotenv) – Leitura do `.env`  

---

## 🚀 Deploy

Pode ser feito em qualquer serviço que suporte Docker ou Go nativamente:

- **Docker Compose**  
- **Heroku**  
- **Render**  
- **Fly.io**  
- **AWS / GCP / Azure**  

Exemplo com **Docker Compose**:

```yaml
version: "3.8"

services:
  api-email:
    build: .
    ports:
      - "8080:8080"
    env_file:
      - .env
```

---

## 📜 Licença

Este projeto está licenciado sob a **MIT License**.  
Você pode usar, modificar e distribuir livremente.

---

## 👨‍💻 Autor

Desenvolvido por [**Marcos Oleniuk**](https://github.com/marcosoleniuk)  
📧 Entre em contato: marcos@moleniuk.com  
