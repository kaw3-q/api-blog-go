# API Blog Backend (Go)

Uma API RESTful para um sistema de blog, desenvolvida em Go, com suporte a autenticação JWT, controle de acesso (RBAC) e persistência em SQLite.

## 🚀 Funcionalidades

- **Autenticação:** Registro e Login de usuários com JWT.
- **Autorização:** Controle de acesso baseado em cargos (User/Admin).
- **Posts:** CRUD completo de postagens.
- **Banco de Dados:** Utiliza SQLite para facilidade de setup.
- **Arquitetura:** Organizada seguindo padrões de Clean Architecture/Hexagonal.

## 🛠️ Tecnologias Utilizadas

- **Go (Golang)**
- **JWT (JSON Web Token)**
- **GORM** (opcional/planejado ou usado via SQL puro conforme implementação)
- **SQLite**

## 📖 Documentação

Para detalhes sobre os endpoints, formatos de requisição e resposta, consulte o arquivo:
👉 [API_DOCS.md](./API_DOCS.md)

## 🚦 Como Executar

1. **Pré-requisitos:** Certifique-se de ter o Go instalado (v1.20+ recomendado).

2. **Clone o repositório:**
   ```bash
   git clone https://github.com/kaw3-q/api-blog-go.git
   cd api-blog-go
   ```

3. **Instale as dependências:**
   ```bash
   go mod tidy
   ```

4. **Execute a aplicação:**
   ```bash
   go run cmd/api/main.go
   ```

A API estará disponível em `http://localhost:8080`.

## 👤 Autor

- **kaw3-q**

---
Desenvolvido como um exemplo de API robusta em Go.
