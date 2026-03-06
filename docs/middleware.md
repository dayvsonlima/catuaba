# Middleware

Every Catuaba app includes a production-ready middleware stack. Middleware lives in the `middleware/` directory and is registered in `application.go`.

## Built-in middleware

The middleware stack runs in order for every request:

| # | Middleware | File | What it does |
|---|-----------|------|-------------|
| 1 | **Health** | `health.go` | `GET /health` returns 200 OK (skips remaining middleware) |
| 2 | **RequestID** | `request_id.go` | Adds `X-Request-ID` header to every response |
| 3 | **Logger** | `logger.go` | Structured JSON logs: method, path, status, latency, request ID |
| 4 | **Recovery** | `recovery.go` | Catches panics, returns 500, logs stack trace |
| 5 | **SecureHeaders** | `secure_headers.go` | Sets HSTS, X-Frame-Options, X-Content-Type-Options, CSP |
| 6 | **CORS** | `cors.go` | Configurable via `CORS_ALLOWED_ORIGIN` env var |
| 7 | **RateLimit** | `rate_limit.go` | Token bucket: 10 requests/second, burst of 40 |
| 8 | **Session** | `session.go` | SCS cookie sessions (OWASP-compliant) |
| 9 | **CSRF** | `csrf.go` | Token validation on POST/PUT/DELETE (skips `/api/` routes) |

Additional helpers (not middleware, but used by handlers):

| Helper | File | Purpose |
|--------|------|---------|
| **Flash** | `flash.go` | `SetFlash(ctx, "success", "Done!")` / `GetFlash(ctx)` |
| **RequireAuth** | `require_auth.go` | Route-level auth guard (only with `--auth`) |

## Registration

Middleware is registered in `application.go`:

```go
engine := gin.New()

engine.Use(middleware.Health())
engine.Use(middleware.RequestID())
engine.Use(middleware.Logger())
engine.Use(middleware.Recovery())
engine.Use(middleware.SecureHeaders())
engine.Use(middleware.CORS())
engine.Use(middleware.RateLimit(10, 40))
engine.Use(middleware.Session())
engine.Use(middleware.CSRF())
```

Order matters. Health comes first so health checks are fast. CSRF comes last because it needs sessions.

## CSRF protection

The CSRF middleware:
- Generates a random token, stored in a cookie (`_csrf_token`)
- Validates the token on `POST`, `PUT`, `DELETE` requests
- Skips `/api/` routes (they use their own auth)
- Forms include a hidden `_csrf_token` field automatically
- HTMX sends the token via `X-CSRF-Token` header (configured in the base layout)

If you need the token in a custom form:

```templ
<form method="POST" action="/my-action">
    <input type="hidden" name="_csrf_token" value={ middleware.CSRFToken(ctx) }/>
    <!-- form fields -->
</form>
```

## Flash messages

Flash messages persist for one request (stored in cookies):

```go
// Set in a handler (before redirect)
middleware.SetFlash(ctx, "success", "Post created!")
middleware.SetFlash(ctx, "error", "Something went wrong")
middleware.SetFlash(ctx, "warning", "Check your input")

// Read in the next handler (the base layout does this automatically)
msg, flashType := middleware.GetFlash(ctx)
```

The base layout renders flash messages with Alpine.js auto-dismiss (5 seconds).

## Auth middleware

When the app is created with `--auth`, a `RequireAuth()` middleware is generated:

```go
// Protect a group of routes
authorized := engine.Group("/admin")
authorized.Use(middleware.RequireAuth())
{
    authorized.GET("/dashboard", admin.Dashboard)
}

// Or protect a single route
engine.GET("/profile", middleware.RequireAuth(), profile.Show)
```

It checks `middleware.GetSessionUserID(ctx)` and redirects to `/login` if not authenticated.

## Session helpers

```go
// Store user in session (after login)
middleware.SetSessionUserID(ctx, user.ID)
middleware.SessionManager.Put(ctx.Request.Context(), "userName", user.Name)

// Read from session
userID := middleware.GetSessionUserID(ctx)
name := middleware.SessionManager.GetString(ctx.Request.Context(), "userName")

// Destroy session (logout)
middleware.SessionManager.Destroy(ctx.Request.Context())
```

## Generating custom middleware

```bash
catuaba g middleware RateLimit
```

Creates `middleware/rate_limit.go` with a stub:

```go
package middleware

import "github.com/gin-gonic/gin"

func RateLimit() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        // TODO: implement rate_limit middleware logic
        ctx.Next()
    }
}
```

Then register it in `application.go`:

```go
engine.Use(middleware.RateLimit())
```

## Configuration

Environment variables used by middleware:

| Variable | Default | Used by |
|----------|---------|---------|
| `CORS_ALLOWED_ORIGIN` | `*` | CORS |
| `SESSION_SECRET` | random | Session |
| `APP_ENV` | `development` | SecureHeaders (HSTS only in production) |

## Learn more

- [Gin middleware](https://gin-gonic.com/docs/examples/custom-middleware/)
- [SCS session manager](https://github.com/alexedwards/scs)
- [OWASP CSRF prevention](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html)
