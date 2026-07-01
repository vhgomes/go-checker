# Go Checker

Serviço de monitoramento de disponibilidade de sites. Cada usuário cadastra URLs, o sistema faz checagens HTTP periódicas em background e expõe status/histórico/dashboard via API REST.

## Funcionalidades

- Autenticação via JWT (registro e login)
- CRUD de sites monitorados, isolado por usuário
- Monitoramento em background com checagem HTTP real (status code, latência, retry com backoff)
- Registro dinâmico de sites recém-criados no monitor, sem necessidade de restart
- Histórico de status por site, com paginação
- Dashboard agregado por usuário, com cache em Redis
- Cronjobs agendados (ex: atualização periódica do dashboard)

## Stack

- **Go** 1.25
- **Gin** — HTTP framework
- **GORM** — acesso a dados (migração para `sqlc + pgx` planejada, ver [Melhorias em andamento](#melhorias-em-andamento))
- **SQLite** — banco de dados
- **Redis** — cache do dashboard
- **JWT** — autenticação
- **Zap** — logging estruturado
- **go-playground/validator** — validação de input

## Arquitetura

```
cmd/
  main.go              # bootstrap: config, DB, redis, router, graceful shutdown
internal/
  config/              # carregamento de configuração via env vars, JWT secret
  handlers/            # camada HTTP (Gin), sem lógica de negócio
  repository/          # acesso a dados, todos os métodos recebem context.Context
  monitor/             # monitoramento em background dos sites
  cronjobs/            # jobs agendados (CronJobManager + implementações)
  server/               # setup de rotas e middlewares
```

Pontos de design relevantes:

- **`context.Context` por requisição**: handlers usam `c.Request.Context()`, nunca guardam contexto em struct. Repositórios recebem `ctx` como primeiro parâmetro e propagam via `.WithContext(ctx)`.
- **`user_id` tipado no contexto**: middleware de autenticação extrai e valida o claim do JWT uma única vez, injetando `uint` já convertido no contexto do Gin. Handlers usam um helper único (`GetUserID(c)`) em vez de repetir type assertions.
- **Monitoramento dinâmico**: novos sites criados via API são registrados no monitor em tempo real (sem restart), através de um mecanismo de registro dinâmico análogo ao `CronJobManager`.

## Como rodar

### Pré-requisitos

- Go 1.25+
- Redis (ex: via `docker-compose up`)

### Variáveis de ambiente

Crie um `.env` na raiz (veja `.env.example`):

| Variável | Descrição | Obrigatória |
|---|---|---|
| `JWT_SECRET` | Chave de assinatura dos tokens JWT | Sim — aplicação falha no boot se ausente |
| `DB_PATH` | Caminho do arquivo SQLite | Não (default local) |
| `REDIS_ADDR` | Endereço do Redis (`host:port`) | Não (default `localhost:6379`) |
| `SERVER_PORT` | Porta HTTP da aplicação | Não (default `8080`) |

### Rodando localmente

```bash
docker-compose up -d          # sobe o Redis
cp .env.example .env          # ajuste as variáveis, principalmente JWT_SECRET
go run cmd/main.go
```

## Endpoints principais

| Método | Rota | Descrição | Autenticado |
|---|---|---|---|
| POST | `/auth/register` | Cria usuário | Não |
| POST | `/auth/login` | Autentica e retorna JWT | Não |
| GET | `/sites` | Lista sites do usuário autenticado | Sim |
| POST | `/sites` | Cria site e inicia monitoramento | Sim |
| PUT | `/sites/:id` | Atualiza site | Sim |
| DELETE | `/sites/:id` | Remove site | Sim |
| GET | `/sites/:id/status` | Histórico de status do site (paginado) | Sim |
| GET | `/dashboard` | Dashboard agregado do usuário (cache Redis) | Sim |

## Segurança

- Senhas com hash via bcrypt
- Mensagens de erro genéricas em login/registro (sem enumeração de usuários)
- Rate limiting em `/auth/login` e `/auth/register`
- Validação de input em todos os payloads (URL, e-mail, intervalos de checagem)
- Toda query de leitura de sites é escopada por `user_id` do token autenticado

## Melhorias em andamento

Itens ainda não concluídos no plano de refatoração:

- [ ] Renomear import aliasing (`handlers2`, `repository2`) causado por colisão de nome de variável com pacote (`redis`)
- [ ] Trocar `panic` em `InitDB` por retorno de erro (`(*gorm.DB, error)`) tratado explicitamente em `main()`
- [ ] Dockerfile multi-stage da aplicação (Go pinado, sem `latest`) + endpoints `/health` e `/ready`
- [ ] Migração de GORM/SQLite para `sqlc + pgx` sobre Postgres
- [ ] Cobertura de testes (unitários com testify, integração com testcontainers)

## Licença

MIT
