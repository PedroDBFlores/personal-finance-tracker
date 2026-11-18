# Personal Finance Tracker — AGENTS.md

Personal finance tracker, tracked as **two projects** in this repo:

```
api/                 Go REST API (Clean Architecture, SQLite via GORM, chi router)
webapp/              React webapp — consumes the api over HTTP
```

- `api/` is a self-contained Go module (`api/go.mod`). All Go commands run from `api/`.
- `webapp/` is a separate Vite + React + TypeScript project with its own tooling; it must
  only talk to the API over HTTP (default `http://localhost:8080`, override with
  `VITE_API_URL`), never import Go code.
- **webapp uses pnpm** (pinned via `"packageManager": "pnpm@12.21.0"` in `webapp/package.json`).
  Use `pnpm install` / `pnpm add` / `pnpm test` / `pnpm build` / `pnpm dev` from `webapp/` —
  never npm/yarn there. `pnpm-lock.yaml` is committed; no `package-lock.json`.
- Shared repo files: `README.md`, this `AGENTS.md`, `.gitignore`.

## webapp/ — structure and commands

```
webapp/src/api.ts        typed API client (fetch wrapper, ApiError with status, base URL from VITE_API_URL)
webapp/src/App.tsx      AppRoutes (the <Routes> table) + default App (BrowserRouter + layout)
webapp/src/pages/        one page component per route (AccountsPage, NewAccountPage, AccountDetailPage, NewMovementPage)
webapp/src/components/RouterLink.tsx  chakra(Link) wrapper: router links with Chakra style props
webapp/src/components/ui/provider.tsx  ChakraProvider wrapper (defaultSystem)
webapp/src/*.test.ts(x)  vitest tests (jsdom); api client tests mock fetch, component tests mock ./api
```

- **Routing: react-router v8.** `App.tsx` exports `AppRoutes` (the `<Routes>` table) so tests
  can render it inside a `MemoryRouter`; the default `App` wraps it in `BrowserRouter`.
  Routes: `/` accounts list, `/accounts/new` create account, `/accounts/:id` account detail
  with movements and investments, `/accounts/:id/movements/new` create movement,
  `/accounts/:id/investments/new` create investment,
  `/accounts/:id/investments/:investmentId` investment detail with transactions,
  `/accounts/:id/investments/:investmentId/transactions/new` create transaction.
  Navigate with `useNavigate`/`<RouterLink to=...>`; page params via `useParams`.
- **UI kit: Chakra UI v3** (`@chakra-ui/react` + `@emotion/react`). Wrap the app in
  `Provider` (`src/components/ui/provider.tsx`) in `main.tsx`; use Chakra components
  (Container, Card, Button, Field, Input, NativeSelect, Badge, ...) with colorPalette /
  variant props for styling. Component tests must render inside `<Provider>` too.
  For router navigation use `RouterLink` (`src/components/RouterLink.tsx`, built with
  `chakra(Link)`) — do NOT use `Button as={Link}`: Chakra's `as` prop does not merge the
  router's props in TypeScript and leaks `to` onto the DOM element.
- Tests: `pnpm test` (vitest run), `pnpm test:watch`. TDD applies here too: failing test first.
- Component tests mock the `./api` module with `vi.mock` and assert on the mocked functions;
  API client tests stub `fetch` with `vi.stubGlobal`.
- Build/lint: `pnpm build` (tsc + vite build), `pnpm lint` (oxlint). Dev server: `pnpm dev` (:5173).
- The API sends permissive CORS headers (`Access-Control-Allow-Origin: *`) so the dev server
  can call it cross-origin.

## api/ — Architecture

Clean Architecture (hexagonal). Dependencies point inward only:
`cmd` -> `adapters` -> `usecases` -> `domain`. `domain` must not import anything from
the rest of the project.

```
api/cmd/             composition root: opens SQLite, builds repos + usecases, serves HTTP on :8080
api/adapters/database/  SQLite persistence: GORM models + repositories implementing usecase ports
api/adapters/rest/      HTTP handlers (chi): one file per entity + router.go
api/usecases/           ports (interfaces) at package root, implementations in usecases/<entity>/,
                        generated fakes in usecases/usecasesfakes/
api/domain/             entities, value objects, enums — no project-internal dependencies
api/dtos/               input/output DTOs used by usecases
api/app_errors/         shared error types (e.g. ResourceNotFoundError)
```

### api/domain/

- One file per concept (`account.go`, `movement.go`, ...). Entities are plain structs
  with a doc comment on every field and a `NewX(...)` constructor.
- Value objects and enums live here too (`Money`, `Id`, `Currency`, `AccountType`,
  `MovementType`, ...). Enums are string types with a const block and a `ParseX(s string)` func.
- IDs: use `domain.Id` (uint, SQLite autoincrement) for all new entities.
  `Investment`/`InvestmentTransaction` still use `uuid.UUID` for their own ids from earlier
  iterations — do not extend that pattern to new code. But `Investment.AccountId` is a
  `domain.Id` FK to `Account` (adjusted when investments were wired up).
- `Money` is `int64` minor units: the last two digits are always the cents.
  `domain.ParseMoney("12.34")` parses a decimal string into minor units (used by the API).
- Investments: `Investment.CurrentValue` is derived, not stored user input — the
  `Get`/`List` investment usecases fetch the investment's transactions and call
  `inv.UpdateCurrentValue(txs)` before returning the DTO. `InvestmentTransaction.Quantity`
  is also `Money`-style (hundredths of a share), so ETF value = `Price * Quantity / 100`.

### api/usecases/

- Ports are interfaces in the root `usecases` package, one file per entity
  (`account_ports.go`, `movement_ports.go`). Each interface carries a
  `//counterfeiter:generate . XxxYyy` directive. Naming: `<Entity>Creator`,
  `<Entity>Updater`, `<Entity>Getter`, `<Entity>Lister`, `<Entity>Deleter`.
  Methods take `ctx context.Context` first. Child entities are scoped by their parent id
  (e.g. movement ports take `accountId`). Updates take a `<Entity>Patch` struct with
  pointer fields (nil = leave unchanged) — partial updates.
- Implementations live in `usecases/<entity>/`, one usecase per file:
  `<verb>_<entity>_usecase.go`, package `<entity>`. Struct `XxxUseCase` holds the ports
  (private fields) plus a `*slog.Logger`; constructor `NewXxxUseCase(ports..., logger)`;
  single entrypoint `Execute(ctx context.Context, ...)`. Map between DTOs and domain
  models inside the usecase; log port calls with `logger.DebugContext(ctx, ...)`.
- Fakes for the ports are generated by counterfeiter into `usecases/usecasesfakes/`
  (`//go:generate go tool counterfeiter -generate` in `usecases/doc.go`). Never edit
  them by hand — run `go generate ./...` after changing a port.

### api/adapters/database/

- GORM models in `models.go` (embed `gorm.Model`), SQLite via `gorm.io/driver/sqlite`.
  Exception: `Investment`/`InvestmentTransaction` use a `uuid.UUID` primary key
  (`gorm:"type:char(36);primaryKey"`) so they declare `CreatedAt`/`UpdatedAt`/`DeletedAt`
  explicitly instead of embedding `gorm.Model`.
- One repository per entity (`AccountRepository`, `MovementRepository`), constructed
  with `NewXxxRepository(db *gorm.DB)`, methods matching the ports exactly.
- Convert between layers with unexported helpers `dbXToDomainX` / `domainXToDBX`
  in the repository file.
- On `Get`/`Update`/`Delete` of a missing row return
  `app_errors.NewResourceNotFoundError("<Entity>", id.String())` (check the error from
  `First` or rows affected == 0). Partial updates use `clause.Assignments` with only the
  non-nil patch fields.

### api/adapters/rest/ (HTTP handlers)

- chi router (`github.com/go-chi/chi/v5`). One handler per entity
  (`account_handler.go`, `movement_handler.go`) plus `router.go` with `NewRouter`.
- Routes: `POST/GET /accounts`, `GET/PATCH/DELETE /accounts/{id}`, nested
  `POST/GET /accounts/{accountId}/movements`, `GET/PATCH/DELETE .../movements/{movementId}`,
  nested `POST/GET /accounts/{accountId}/investments`,
  `GET/PATCH/DELETE .../investments/{investmentId}` (uuid param, parsed with `uuid.Parse`),
  and nested `.../investments/{investmentId}/transactions[/{transactionId}]`.
- Handlers depend on narrow interfaces matching the usecase `Execute` signatures
  (satisfied by the usecase structs). Request/response are JSON wire types defined in
  the handler file; money crosses the wire as a decimal string (`"25.99"`), converted
  with `domain.ParseMoney` / `Money.String()`. Empty money fields in transaction
  create requests mean zero (`parseMoneyParam`).
- Error mapping: `ResourceNotFoundError` -> 404, invalid body/params -> 400,
  anything else -> 500; error responses are `{"error": "..."}`.
- **Handler tests use the same pattern as usecase tests**: build the real usecase with a
  counterfeiter fake port (`account.NewCreateAccountUseCase(&usecasesfakes.FakeAccountCreator{}, logger)`),
  inject it into the handler, drive requests with `httptest`, and assert on the port via
  `XxxArgsForCall(0)` / `XxxReturns(...)`. Do not hand-write fake usecases.

### api/dtos/ and api/app_errors/

- DTOs are plain structs with `NewXxxDto(...)` constructors; read DTOs like
  `AccountDto` embed the domain entity.
- Application errors live in `app_errors` with `New...Error` constructors and an
  `Error()` method. Usecases propagate port errors unchanged (tests assert with
  `require.ErrorIs`).

## Testing — strict TDD (red -> green -> refactor)

1. Write the failing test first.
2. Implement the minimum to make it pass.
3. Refactor.

Every entity, usecase, repository and handler ships with tests; a feature is not done
without them.

- Test files use the external test package (`package account_test`, `package rest_test`, ...),
  `testing` + `github.com/stretchr/testify/require`, and `t.Run` subtests per scenario
  (happy path + error path).
- Usecase and handler tests fake the ports with `usecasesfakes.FakeXxx`, set behavior with
  `XxxReturns(...)`, and assert calls with `XxxArgsForCall(0)`. Logger:
  `slog.New(slog.NewTextHandler(os.Stdout, nil))` (or `io.Discard` in handler tests); ctx: `context.TODO()`.
- Repository tests run against in-memory SQLite using the `getInMemoryDatabase(t, logger)`
  helper in `adapters/database/database_test.go` (`sqlite.Open(":memory:")` + `AutoMigrate`).
- Domain tests are plain unit tests.

## Commands (run from `api/`)

Verified with Go 1.27:

```
go build ./...
go test ./...
go vet ./...
go generate ./...   # regenerate counterfeiter fakes after changing ports
go run ./cmd        # serve the API on :8080 (SQLite file: tracker.db)
```

## Housekeeping

- Schema changes happen in `api/adapters/database/models.go` (AutoMigrate). `SCHEMA.md` still
  describes the old Badger KV layout — update it when the schema changes.
- `github.com/dgraph-io/badger` in `api/go.mod` and `api/adapters/database/connector.go` are legacy
  (pre-SQLite); the composition root in `api/cmd` opens SQLite directly via
  `gorm.Open(sqlite.Open("tracker.db"), ...)` + `AutoMigrate`. Do not use Badger in new code.
- `*.db` files are gitignored.
- Git: work on `main`, short imperative commit messages.
