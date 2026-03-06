# Catuaba

A Rails-like CLI framework for Go. Generates full-stack web applications with Gin, GORM, Templ, HTMX, and Tailwind CSS.

## Quick Reference

```bash
go build -v .                    # Build the CLI
go test -v ./...                 # Run all tests
go test -v ./cli/new/...         # Test specific package
go test -v -run TestAction ./cli/scaffold/...  # Run specific test
```

## Architecture

```
catuaba.go              # CLI entrypoint — all commands registered here (urfave/cli)
cli/                    # Command implementations — one package per command
  new/action.go         # `catuaba new <name>` — creates full app
  scaffold/action.go    # `catuaba g scaffold` — HTML handlers + views + routes
  api/action.go         # `catuaba g api` — JSON API controllers + routes
  model/action.go       # `catuaba g model` — model + SQL migration
  controller/action.go  # `catuaba g controller`
  view/action.go        # `catuaba g view`
  migration/action.go   # `catuaba g migration`
  middleware/action.go   # `catuaba g middleware`
  service/action.go     # `catuaba g service`
  seed/action.go        # `catuaba g seed`
  db/migrate.go         # `catuaba db migrate|rollback|status`
  routes/action.go      # `catuaba routes` — lists app routes
  server/action.go      # `catuaba server`
  install/action.go     # `catuaba install` — plugin installer
  mcp/action.go         # `catuaba mcp` — starts MCP server
  output/output.go      # Colored CLI output helpers
generator/              # Template engine
  embed.go              # `//go:embed all:templates` — embeds all .tmpl files
  render.go             # Render() / RenderFromContent() — Go text/template
  generate_file.go      # GenerateFile() / GenerateFromContent() — write to disk
  mkdir.go              # Mkdir() — create directories
  helpers.go            # Template functions (toModelName, toSnake, toPlural, etc.)
  module.go             # ModuleName() — reads module name from go.mod in cwd
  project.go            # IsInsideCatuabaProject() — validates cwd is a Catuaba app
  templates/            # All .tmpl files (embedded at compile time)
    application/        # Templates for `catuaba new` (full app skeleton)
    scaffold/           # Templates for `catuaba g scaffold` and `catuaba g api`
    model.go.tmpl       # Standalone model template
    controller.go.tmpl  # Standalone controller template
    *.tmpl              # Other generator templates
code_editor/            # AST-based Go source code manipulation
  add_import.go         # AddImport() / AddAliasedImport() — inject imports via AST
  insert_attribute.go   # InsertAttribute() — add fields to structs
  edit_file.go          # EditFile() — read-modify-write pattern
mcp/                    # MCP server for AI integration
  server.go             # Server setup
  parser/               # Go source code parsers (routes, models, controllers, etc.)
plugin/                 # Plugin system (manifest, registry, install, inject)
```

## Key Patterns

### Template System

Templates use Go `text/template` syntax and live in `generator/templates/`. They are embedded via `embed.FS` at compile time.

**File naming convention:**
- `.go.tmpl` — generates a `.go` file
- `.templ.tmpl` — generates a `.templ` file (Templ views). The outer `.tmpl` is consumed by Go template; the inner `.templ` is the output extension
- `.sql.tmpl` — generates SQL migration files

**Template data structs:**

| Command | Data struct | Key fields |
|---------|-------------|------------|
| `catuaba new` | `new.AppBuilder` | `.Name`, `.DBDriver`, `.DBHost`, `.DBPort`, `.DBUser`, `.DBName`, `.Auth` |
| `catuaba g model/scaffold/api` | `model.ModelBuilder` | `.Name` (CamelCase), `.Params` (slice of `"name:type"`) |
| scaffold handlers | `scaffold.HandlerBuilder` | `.Name`, `.MethodName`, `.Params` |
| api controllers | `api.ControllerBuilder` | `.Name`, `.MethodName`, `.Params` |

**Template functions** (available in all templates via `templateFuncMap()`):

| Function | Example | Result |
|----------|---------|--------|
| `toModelName` / `camelize` | `post_comment` | `PostComment` |
| `toSnake` | `PostComment` | `post_comment` |
| `toPlural` | `Post` | `Posts` |
| `toLowerPlural` | `Post` | `posts` |
| `toVarName` | `post_comment` | `postComment` |
| `toAttrName` | `title:string` | `Title` |
| `toType` | `title:string` → `string`, `count:integer` → `int`, `active:boolean` → `bool` |
| `toRawType` | `title:string` | `string` (as-is) |
| `toJson` | `title:string` | `` `json:"title"` `` |
| `toJsonBinding` | `title:string` | `` `json:"title" binding:"required"` `` |
| `toFormBinding` | `title:string` | `` `form:"title"` `` |
| `toFormLabel` | `PostId` | `Post Id` |
| `toFormInput` | `text` → `textarea`, `boolean` → `checkbox`, `integer` → `number` |
| `toSQLType` | `string` → `VARCHAR(255)`, `text` → `TEXT`, `integer` → `INTEGER` |
| `toSQLDefault` | `string` → `''`, `integer` → `0`, `boolean` → `0` |
| `toColumnName` | `title:string` | `title` (snake_case) |
| `moduleName` | (zero-arg) | reads module name from cwd's `go.mod` |

**Critical distinction:**
- `{{.Name}}` in `application/` templates = app name (raw, as provided by user)
- `{{.Name}}` in `scaffold/` templates = CamelCase resource name (e.g., `Post`)
- `{{moduleName}}` = reads actual `go.mod` at generation time (use in scaffold/api templates for imports)

### Generation Flow

1. `GenerateFile(tmplName, data, destPath)` — loads embedded template, renders with data, writes to `cwd/destPath`
2. `GenerateFromContent(content, data, destPath)` — same but from raw string (used for `.keep` files)
3. All paths are relative to `os.Getwd()`

### Route Injection

Scaffold and API generators inject routes into `config/routes.go` using:
1. `code_editor.EditFile()` — read-modify-write
2. `code_editor.AddImport()` / `AddAliasedImport()` — AST-based import injection
3. Routes appended before closing `}` of `SetupRoutes()`

API routes use aliased imports (e.g., `api_posts "myapp/app/controllers/api/posts"`) to avoid collisions with scaffold handler imports.

### Scaffold vs API

| | `catuaba g scaffold` | `catuaba g api` |
|---|---|---|
| Handlers | `app/controllers/{resource}/` (HTML) | `app/controllers/api/{resource}/` (JSON) |
| Views | `app/views/{resource}/*.templ` | None |
| Routes | `GET /posts`, `POST /posts`, etc. | `GET /api/posts`, `POST /api/posts`, etc. |
| Route template | `scaffold/view_routes.go.tmpl` | `scaffold/api_routes.go.tmpl` |

### Generated App Stack

| Component | Technology |
|-----------|-----------|
| HTTP framework | Gin |
| ORM | GORM (postgres, mysql, sqlite) |
| Views | Templ (type-safe, compiled) |
| Interactivity | HTMX + Alpine.js |
| CSS | Tailwind CSS (standalone CLI) |
| Sessions | SCS |
| Auth | Optional (`--auth` flag) — bcrypt via `golang.org/x/crypto` |
| DB migrations | Raw SQL files (up/down) via golang-migrate |

### Conditional Auth

`catuaba new myapp --auth` generates extra files:
- `app/models/user.go` — User model with password hash
- `app/controllers/auth/` — login, register, logout handlers
- `app/views/auth/` — login, register Templ views
- `middleware/require_auth.go` — session-based auth middleware
- `database/migrations/00000000000001_create_users.{up,down}.sql`

Templates use `{{if .Auth}}...{{end}}` for conditional blocks.

## Testing Conventions

- Use `testify/assert` and `testify/require`
- Test files colocated: `action.go` → `action_test.go`
- For CLI actions: use `t.TempDir()`, `os.Chdir(tmpDir)`, create a `cli.Context` with flags
- Check generated file existence with `assert.FileExists` / `assert.DirExists`
- Check generated content with `os.ReadFile` + `assert.Contains`
- Always `defer os.Chdir(origDir)` after changing directories

## Adding a New Generator

1. Create template(s) in `generator/templates/`
2. Create CLI action in `cli/<name>/action.go`
3. Register command in `catuaba.go` under the `generate` subcommands
4. If it modifies `config/routes.go`, use `code_editor.EditFile()` + `code_editor.AddImport()`
5. Add tests in `cli/<name>/action_test.go`

## Adding Templates to `catuaba new`

1. Create `.tmpl` file in `generator/templates/application/`
2. Add directory to `dirs` slice in `cli/new/action.go` if needed
3. Add `{tmpl, dest}` pair to `files` slice in `cli/new/action.go`
4. Update test in `cli/new/action_test.go`
