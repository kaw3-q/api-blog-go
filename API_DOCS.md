# Documentação da API - Blog Backend (Go)

Esta API fornece um sistema completo de blog com autenticação JWT, controle de acesso baseado em cargos (RBAC) e persistência em banco de dados.

## Informações Gerais
- **Base URL:** `http://localhost:8080`
- **Formato de Dados:** `application/json`
- **Autenticação:** Bearer Token (JWT)

---

## 1. Autenticação e Usuários

### Registrar Usuário
Cria uma nova conta no sistema.
- **URL:** `/register`
- **Método:** `POST`
- **Corpo da Requisição:**
```json
{
  "username": "joao_dev",
  "email": "joao@email.com",
  "password": "senha_segura",
  "role": "user" 
}
```
> Nota: O campo `role` pode ser `user` ou `admin`. Se omitido, o padrão é `user`.

### Login
Autentica um usuário e retorna um token JWT.
- **URL:** `/login`
- **Método:** `POST`
- **Corpo da Requisição:**
```json
{
  "email": "joao@email.com",
  "password": "senha_segura"
}
```
- **Resposta de Sucesso (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "joao_dev",
    "email": "joao@email.com",
    "role": "user",
    "created_at": "2026-06-13T10:00:00Z"
  }
}
```

---

## 2. Postagens (Blog)

### Listar Posts
Retorna todas as postagens cadastradas.
- **URL:** `/posts`
- **Método:** `GET`
- **Autenticação:** Não requerida (Pública).

### Criar Post
Adiciona uma nova postagem ao blog.
- **URL:** `/posts`
- **Método:** `POST`
- **Autenticação:** Requerida (Bearer Token).
- **Cabeçalho:** `Authorization: Bearer <TOKEN>`
- **Corpo da Requisição:**
```json
{
  "title": "Minha Primeira Postagem",
  "content": "Conteúdo detalhado do post aqui...",
  "author_id": 1
}
```

---

## 3. Administração (Acesso Restrito)

### Listar Todos os Usuários
Retorna a lista de todos os usuários registrados.
- **URL:** `/admin/users`
- **Método:** `GET`
- **Autenticação:** Requerida (Token deve pertencer a um usuário com role `admin`).
- **Cabeçalho:** `Authorization: Bearer <TOKEN>`

---

## Códigos de Erro Comuns

| Código | Descrição | Motivo Comum |
| :--- | :--- | :--- |
| **400** | Bad Request | JSON malformado ou campos obrigatórios ausentes. |
| **401** | Unauthorized | Token ausente, expirado ou inválido. |
| **403** | Forbidden | Usuário autenticado, mas sem permissão (ex: não é admin). |
| **404** | Not Found | Recurso (post ou usuário) não encontrado. |
| **405** | Method Not Allowed | Uso de GET em rota que aceita apenas POST, etc. |

---

## Dicas para o Front-end
1. **Persistência do Token:** Armazene o `token` retornado no login no `localStorage` ou em um `Cookie` seguro.
2. **Interceptação:** Configure seu cliente HTTP (Axios/Fetch) para incluir o cabeçalho `Authorization: Bearer <seu_token>` em todas as requisições para rotas protegidas.
3. **Segurança:** Nunca armazene a senha do usuário no estado global da aplicação.
