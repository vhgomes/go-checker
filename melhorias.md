# Go Checker — Checklist de Refatoração

Plano de ação para elevar o projeto a nível sênior. Ordem sugerida: resolva todos os `[CRÍTICO]` antes de qualquer outra coisa — são bugs ativos em produção, não dívida técnica estética.

---

## 🔴 CRÍTICO — Bugs ativos / segurança

- [x] **IDOR em `GetSitesByUserId`** (`internal/repository/site_repo.go`)
  `r.DB.Find(&sites).Where("user_id = ?", userId)` — `.Find()` executa antes do `.Where()` ser aplicado. `GET /sites` retorna sites de **todos os usuários**, não só do autenticado.
  **Fix:** inverter ordem → `Where(...).Find(&sites)`.

- [x] **JWT secret hardcoded** (`internal/config/jwt.go`)
  `var JwtSecret = []byte("minha_chave_super_secreta")` versionado no git — qualquer um forja tokens válidos.
  **Fix:** carregar de env var no boot; `log.Fatal` se ausente.

- [x] **Chave do Redis corrompida no cronjob** (`internal/cronjobs/cronjob_dashboard.go:34`)
  `string(rune(user))` converte o ID em code point Unicode, não em texto. Cache do dashboard grava/lê chaves erradas.
  **Fix:** `fmt.Sprintf("dashboard:user:%d", user)` ou `strconv.Itoa`.

- [x] **Novos sites nunca entram em monitoramento**
  `StartMonitoring` roda uma única vez no boot. `CreateSite` só faz `INSERT`, não dispara monitoramento.
  **Fix:** mecanismo de registro dinâmico — canal de novos sites ou um `MonitorManager` com polling/diff, análogo ao `CronJobManager` existente.

- [x] **Monitoramento é fake** (`internal/monitor/monitor.go` — `checkSiteRandom`)
  Status e response time são sorteados com `math/rand`, nenhuma chamada HTTP real acontece.
  **Fix:** `net/http` com `http.Client{Timeout}`, request real usando o `checkCtx` já existente, retry com backoff, captura de status/latência reais.

---

## 🟠 IMPORTANTE — Robustez, segurança secundária, consistência com o stack

- [x] **Type assertion sem verificação repetida em ~10 handlers**
  `userAny, _ := c.Get("user_id"); userID := uint(userAny.(float64))` — panic se o tipo vier diferente.
  **Fix:** middleware injeta `userID uint` já validado no contexto; handlers usam helper único (`GetUserID(c)`).

- [x] **User enumeration no login** (`internal/repository/user_repo.go` — `Login`)
  Retorna `"usuário não encontrado"` vs `"senha incorreta"` como erros distintos, propagados ao cliente.
  **Fix:** mensagem genérica única ("credenciais inválidas") nos dois casos; detalhe só em log.

- [x] **`context.Context` guardado como campo de struct** (`SiteHandler.ctx`)
  Vive pra sempre em vez de vir de `c.Request.Context()` a cada request — cancelamento por request não se propaga.
  **Fix:** usar `c.Request.Context()` dentro de cada handler.

- [x] **Repos sem `context.Context` na maioria dos métodos**
  `AddSite`, `UpdateSite`, `DeleteSite`, `GetSiteById`, `GetSitesByUserId` e todo `SiteStatusRepo` não recebem/propagam `ctx`.
  **Fix:** todo método de repo que toca o banco recebe `ctx` como primeiro parâmetro e usa `.WithContext(ctx)`.

- [x] **Race condition em `CreateUser`** (check-then-create)
  Duas requests concorrentes com o mesmo email podem passar pelo check antes de qualquer insert.
  **Fix:** remover check prévio, inserir direto, tratar erro de unique constraint especificamente.

- [x] **Zero validação de input**
  URL sem validar formato/scheme, `check_interval` sem limite mínimo, email/senha sem validação em `RegisterUser`.
  **Fix:** `go-playground/validator` com tags de binding + rate limiter básico em `/login` e `/register`.

---

## 🟡 MELHORIA — Qualidade, manutenibilidade, prontidão para produção

- [x] **Sem logging estruturado** — trocar `log.Println`/`log.Printf` por Zap com campos (`user_id`, `site_id`, `err`).
- [x] **Configuração hardcoded** — `localhost:6379`, `test.db`, porta `:8080` fixos no código. Criar `internal/config` carregado de env vars.
- [ ] **Naming de import colidindo com pacote** (`handlers2`, `repository2` em `server.go`) — renomear variável `redis` para `redisClient` e remover os aliases forçados.
- [ ] **`panic` em `InitDB`** fora de `init()`/`main()` — trocar por `(*gorm.DB, error)` retornado e `log.Fatal` explícito em `main()`.
- [ ] **Sem health check nem Dockerfile da aplicação** — adicionar `Dockerfile` multi-stage com Go pinado (`golang:1.25-alpine3.19`, nunca `latest`) e endpoints `/health`/`/ready`.
- [x] **Rota `GetDashboardByUser` órfã** — handler existe mas não está registrado em nenhuma rota; registrar ou remover.

---

## Ordem de ataque sugerida

1. IDOR (`GetSitesByUserId`) e JWT secret — maior peso em qualquer avaliação de segurança/senioridade.
2. Extração do middleware de `userID` — resolve o IMPORTANTE mais repetido de uma vez só.
3. Monitoramento real via HTTP + registro dinâmico de novos sites — é a feature principal do produto.
4. Migração para sqlc + Postgres — trabalho estrutural maior, fazer incremental por repo.
5. Restante dos itens `MELHORIA` conforme tempo disponível.
