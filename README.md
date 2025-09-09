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

## 📂 Estrutura do Projeto
api-go-send-email/
│── config/ # Configurações (carregamento de env, smtp, etc.)
│── database/ # (Opcional) Persistência de dados ou logs
│── handlers/ # Handlers HTTP com as rotas
│── middleware/ # Middlewares (ex: autenticação, logs, CORS)
│── models/ # Definições de modelos (request/response)
│── .dockerignore # Arquivos ignorados no build Docker
│── .env.example # Exemplo de variáveis de ambiente
│── .gitignore # Arquivos ignorados no Git
│── Dockerfile # Containerização da aplicação
│── go.mod # Gerenciamento de dependências
│── go.sum # Checksum das dependências
│── main.go # Ponto de entrada da aplicação

## ⚙️ Configuração

### Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto, baseado no `.env.example`:

```env
# Configuração do servidor SMTP
SMTP_HOST=smtp.seuprovedor.com
SMTP_PORT=587
SMTP_USER=seuusuario
SMTP_PASS=suasenha

# Email remetente
SENDER_EMAIL=seuemail@dominio.com

# Porta da API
APP_PORT=8080
```

# Clone o repositório
git clone https://github.com/marcosoleniuk/api-go-send-email.git
cd api-go-send-email

# Instale dependências
go mod tidy

# Execute a aplicação
go run main.go
