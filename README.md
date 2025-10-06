# Go Checker

Go Checker é uma aplicação de monitoramento de sites desenvolvida em Go, com autenticação via JWT, dashboard de status de sites, histórico de monitoramento e armazenamento de informações em SQLite e Redis.

---

## Funcionalidades

- Registro e login de usuários com JWT
- Criação, atualização e exclusão de sites
- Monitoramento contínuo de sites (status e tempo de resposta)
- Histórico de status de sites com filtros por data e status
- Dashboard do usuário com métricas gerais e últimos eventos
- Armazenamento temporário de dashboards em Redis para rápida consulta
- Cronjob para atualização periódica do dashboard

---

## Tecnologias utilizadas

- [Go](https://golang.org/)
- [Gin](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [SQLite](https://www.sqlite.org/index.html)
- [Redis](https://redis.io/)
- [JWT](https://github.com/golang-jwt/jwt)
- [Cron](https://github.com/robfig/cron)

---

## Estrutura do Projeto
```markdown
.
├── cmd/                # Arquivo main.go
├── config/             # Configurações de banco e JWT
├── cronjobs/           # Jobs periódicos (dashboard)
├── handlers/           # Handlers de API
├── monitor/            # Monitoramento dos sites
├── repository/         # Repositórios e entidades
├── server/             # Setup do servidor e rotas
└── utils/              # Funções utilitárias (JWT, etc)

```

---

## Como rodar o projeto

1. Clone o repositório:

```bash
git clone https://github.com/vhgomes/go-checker.git
cd go-checker
````

2. Instale as dependências:

```bash
go mod tidy
```

3. Execute o projeto:

```bash
go run cmd/main.go
```

O servidor estará rodando em `http://localhost:8080`.

---

## Rotas principais

### Usuário

* `POST /register` → Registrar novo usuário
* `POST /login` → Login e obtenção do token JWT

### Sites (autenticado)

* `POST /sites` → Criar site
* `GET /sites` → Listar todos os sites do usuário
* `GET /sites/:id` → Detalhes de um site
* `PUT /sites/:id` → Atualizar site
* `DELETE /sites/:id` → Deletar site

### Status do site (autenticado)

* `GET /sites/:id/status` → Histórico paginado
* `GET /sites/:id/status/date` → Histórico filtrado por datas
* `GET /sites/:id/status/filter` → Histórico filtrado por status
* `GET /sites/:id/status/last` → Último status
* `GET /sites/:id/status/first` → Primeiro status

---

## Mudanças futuras / melhorias (melhorias recomendadas por IA)

* Padronizar campos GORM e usar `gorm.Model` para consistência
* Passar repositórios como ponteiros nos handlers para evitar cópias
* Melhorar tratamento de erros com mensagens mais descritivas
* Corrigir conversão de `user_id` do JWT e usar tipo seguro
* Validação mais completa de inputs (URLs, intervalos, campos obrigatórios)
* Reduzir duplicação em queries de banco e filtrar dados de forma genérica
* Garantir uso de `context.Context` em todas operações de banco
* Gerenciamento mais robusto de cronjobs e persistência de estado
* Implementar monitoramento real de sites via HTTP com timeout e retries
* Tornar JWT configurável via variável de ambiente e considerar refresh tokens
* Corrigir uso de `.Find().Where()` no GORM e outras pequenas inconsistências
* Adotar logger estruturado e consistente em todos os pacotes



