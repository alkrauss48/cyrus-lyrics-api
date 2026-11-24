# AGENTS.md

Guidance for AI coding agents (Claude Code, Codex, Cursor, Copilot, etc.) working in this repository. `CLAUDE.md` imports this file, so keep shared instructions here.

## Project at a glance

This is the backend API for the **CyrusLyrics iOS app** (https://github.com/alkrauss48/cyrus-lyrics). It runs in production at `https://api.cyruskrauss.com`. It handles Google OAuth sign-in and lets signed-in users create, list, read, and trash their own Google Sheets of songs.

- Language: Go 1.26+ (Docker and CI build with 1.27), using the [Gin](https://github.com/gin-gonic/gin) web framework.
- Google APIs: Drive v3 (list and trash files) and Sheets v4 (create and read). The OAuth scope is **`drive.file`**, so the app can only see files that it created itself. Don't widen the scope. It would trigger Google's restricted-scope verification.
- The server is stateless: no database and no sessions. The client holds the user's OAuth token and sends it with each request.

## Repository layout

```
main.go           Router, http.Server with timeouts, graceful shutdown (port 8000)
public/           Unauthenticated routes: GET / and GET /sheets/default
oauth/            GET /oauth/google (login redirect) and /oauth/google/callback
sheets/           Authenticated Google Sheets/Drive routes (index, show, create, delete)
helpers/          OAuth config from env vars, token parsing, Google clients, error -> status mapping
docker/local/     Dev image (go run, source mounted by docker-compose)
docker/prod/      Multi-stage production image, built and pushed by CI
.github/          CI (lint, test + govulncheck, then build and push the image) and Dependabot
```

## Contract with the iOS app

The iOS app is a separate repo, and released versions are installed on user devices that you can't update. **Treat every route, parameter, status code, and response shape as a public API.** Don't make breaking changes without asking.

- **Sign-in**: the app opens `/oauth/google`. Google redirects to `/oauth/google/callback`, which redirects to `cyruslyrics://login?token=<token>&...`. The app reads only the `token` query item. The `access_token`, `refresh_token`, and `expiry` params serve pre-2.0 app versions.
- **Token format**: `token` is standard base64 of the JSON-encoded `oauth2.Token`. Clients send it back as `Authorization: Bearer <token>`. The deprecated `?access_token=&refresh_token=&expiry=` query auth is still accepted for pre-2.0 apps.
- **Status codes drive sign-out** (current app): the app **signs the user out on any 4xx** from `GET /sheets` and `GET /sheets/:id`, and on 401/403 from create and delete. It keeps the user signed in on 5xx and network errors. Use 4xx only when the token or request is actually bad. Use 5xx for upstream or server failures. `helpers.StatusForError` applies this policy, so route errors through `helpers.AbortWithError`.
- **Released v2.1 and earlier ignore status codes**: they decode the body as the success type and sign out if decoding fails. An error body must therefore **never decode as a success response**. The success structs in those versions have non-optional fields (`files`; `majorDimension`, `range`, `values`), and `{"error": "..."}` satisfies none of them.
- **Response bodies**: success responses pass through Google's JSON (`{"files":[{"id","name"}]}` for the index, the `ValueRange` for show, and `{"spreadsheetId"}` for create). Errors are `{"error": "<message>"}`. No app version reads error bodies.
- **Sheet format**: row 1 is headers. Columns are `Category | Sub-category | Song name | URL | Lyrics | Spotify URL`. `POST /sheets` writes these headers as `Genre | Artist | Song Title | URL | Lyrics | Spotify URL`, and `GET /sheets/:id` reads `A:F`.

## Configuration

All configuration comes from environment variables. Copy `.env.example` to `.env` for local development. `.env` is gitignored, so never commit real credentials.

- Google OAuth (required): `CLIENT_ID`, `CLIENT_SECRET`, `REDIRECT_URI`. Optional: `AUTH_URI` and `TOKEN_URI`, which default to Google's endpoints. Tests point `TOKEN_URI` at a fake server.
- Default sheets returned by `/sheets/default`: `DEMO_SHEET_ID`, `AARON_SHEET_ID`
- Gin: `GIN_MODE=release` in production

## Building, testing, running

```sh
go build ./...
go test ./...
go vet ./...
golangci-lint run        # golangci-lint v2, config in .golangci.yml
docker compose up        # local server on :8000, needs .env
```

CI (`.github/workflows/ci.yml`) runs golangci-lint, `go test -race`, and `govulncheck`. On every push it then builds `docker/prod/Dockerfile` and pushes the image to Docker Hub. If you change the Go version, update `go.mod`, both Dockerfiles, and `go_version` in CI together.

## Conventions

- Run `gofmt` on all Go code. Use tabs for indentation in Go files.
- Keep handlers thin. Use `helpers.GetDriveService` or `helpers.GetSheetsService`, pass `c.Request.Context()` to Google calls, and `return` right after every error response.
- Report failures with `helpers.AbortWithError` (the status is mapped from the error) or `helpers.AbortWithMessage` (request validation). Validate request params before doing any OAuth work.
- Add tests next to the code (`*_test.go`) using the standard library and `net/http/httptest` against a Gin router. Tests must never call Google. Test the paths that fail before any network call, or point `TOKEN_URI` at an `httptest` server.
- Keep changes small and focused. This is a personal project that serves a released App Store app.

## Gotchas for agents

- Never log or echo OAuth tokens. They travel in the `Authorization` header and in the legacy query params. The request logger skips query strings for this reason, so keep it that way.
- The OAuth `state` param is a fixed string (`"state-token"`) and isn't verified in the callback. This is a known gap. Fixing it changes the sign-in flow, so coordinate with the iOS app.
- The `README.md` route list is user-facing documentation. Update it when routes change.
- `GET /` is a placeholder message. Its content is the maintainer's call, so don't change it unless asked.
