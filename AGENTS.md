# AGENTS.md

Guidance for AI coding agents (Claude Code, Codex, Cursor, Copilot, etc.) working in this repository. `CLAUDE.md` imports this file, so keep shared instructions here.

## Project at a glance

This is the backend API for the **CyrusLyrics iOS app** (https://github.com/alkrauss48/cyrus-lyrics). It runs in production at `https://api.cyruskrauss.com`. It handles Google OAuth sign-in and lets signed-in users create, list, read, and trash their own Google Sheets of songs.

- Language: Go, using the [Gin](https://github.com/gin-gonic/gin) web framework.
- Google APIs: Drive v3 (list and trash files) and Sheets v4 (create and read). The OAuth scope is **`drive.file`**, so the app can only see files that it created itself. Don't widen the scope. It would trigger Google's restricted-scope verification.
- The server is stateless: no database and no sessions. The client holds the user's OAuth token and sends it with each request.

## Repository layout

```
main.go           Router setup and server startup (port 8000)
public/           Unauthenticated routes: GET / and GET /sheets/default
oauth/            GET /oauth/google (login redirect) and /oauth/google/callback
sheets/           Authenticated Google Sheets/Drive routes (index, show, create, delete)
helpers/          OAuth config from env vars, token parsing, authenticated HTTP client
docker/local/     Dev image (go run, source mounted by docker-compose)
docker/prod/      Multi-stage production image, built and pushed by CI
.github/workflows CI: lint, test, then build and push the Docker image
```

## Contract with the iOS app

The iOS app is a separate repo, and released versions are installed on user devices that you can't update. **Treat every route, parameter, status code, and response shape as a public API.** Don't make breaking changes without asking.

- **Sign-in**: the app opens `/oauth/google`. Google redirects to `/oauth/google/callback`, which redirects to `cyruslyrics://login?token=<token>&...`. The app reads only the `token` query item. The `access_token`, `refresh_token`, and `expiry` params serve pre-2.0 app versions.
- **Token format**: `token` is standard base64 of the JSON-encoded `oauth2.Token`. Clients send it back as `Authorization: Bearer <token>`. The deprecated `?access_token=&refresh_token=&expiry=` query auth is still accepted for pre-2.0 apps.
- **Status codes drive sign-out**: the app **signs the user out on any 4xx** from `GET /sheets` and `GET /sheets/:id`, and on 401/403 from create and delete. It keeps the user signed in on 5xx and network errors. Use 4xx only when the token or request is actually bad. Use 5xx for upstream or server failures.
- **Response bodies**: success responses pass through Google's JSON (`{"files":[{"id","name"}]}` for the index, the `ValueRange` for show, and `{"spreadsheetId"}` for create). The app doesn't parse error bodies.
- **Sheet format**: row 1 is headers. Columns are `Category | Sub-category | Song name | URL | Lyrics | Spotify URL`. `POST /sheets` writes these headers as `Genre | Artist | Song Title | URL | Lyrics | Spotify URL`, and `GET /sheets/:id` reads `A:F`.

## Configuration

All configuration comes from environment variables. Copy `.env.example` to `.env` for local development. `.env` is gitignored, so never commit real credentials.

- Google OAuth: `CLIENT_ID`, `CLIENT_SECRET`, `REDIRECT_URI`, `PROJECT_ID`, `AUTH_URI`, `TOKEN_URI`, `AUTH_PROVIDER_CERT_URL`
- Default sheets returned by `/sheets/default`: `DEMO_SHEET_ID`, `AARON_SHEET_ID`
- Gin: `GIN_MODE=release` in production

## Building, testing, running

```sh
go build ./...
go test ./...
go vet ./...
docker-compose up        # local server on :8000, needs .env
```

## Conventions

- Run `gofmt` on all Go code. Use tabs for indentation in Go files.
- Keep handlers thin. Shared Google client and token logic belongs in `helpers/`.
- Add tests next to the code (`*_test.go`) using `net/http/httptest` against a Gin router. Tests must not call Google, so test the paths that fail before any network call.
- Keep changes small and focused. This is a personal project that serves a released App Store app.

## Gotchas for agents

- Never log or echo OAuth tokens. They travel in the `Authorization` header and in the legacy query params.
- The `README.md` route list is user-facing documentation. Update it when routes change.
- `GET /` is a placeholder message. Its content is the maintainer's call, so don't change it unless asked.
