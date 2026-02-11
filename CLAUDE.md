# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Mindwell-server is a social blogging platform written in Go, featuring blog entries, comments, direct messaging, real-time notifications, and image processing. The system uses PostgreSQL with RUM full-text search extension and implements OAuth2 authentication.

## Build, Run, and Test Commands

### Code Generation
```bash
# Generate API server code from Swagger spec
./scripts/generate.sh

# This generates:
# - restapi/ (main API server scaffolding)
# - restapi_images/ (image server scaffolding)
# - models/ (data models from Swagger spec)
```

### Database Setup
```bash
# Create database and schema
psql -c 'create database mindwell'
psql -d mindwell -q -f scripts/mindwell.sql

# Apply updates (when needed)
psql -d mindwell -q -f scripts/update.sql
```

### Configuration
```bash
# Main server configuration
cp configs/server.sample.toml configs/server.toml
nano configs/server.toml

# Image server configuration
cp configs/images.sample.toml configs/images.toml
nano configs/images.toml
```

### Running Tests
```bash
# Run all integration tests
go test ./test/ --failfast

# Run specific test
go test ./test/ -run TestName --failfast

# Run tests with verbose output
go test -v ./test/ --failfast
```

### Running Servers
```bash
# Run main API server (default port 8000)
go run ./cmd/mindwell-server/ --port 8000

# Run image processing server (default port 8888)
go run ./cmd/mindwell-images-server/ --port 8888

# Run helper CLI tool
go run ./cmd/mindwell-helper/
```

### Building
```bash
# Build all binaries
go build ./cmd/mindwell-server/
go build ./cmd/mindwell-images-server/
go build ./cmd/mindwell-helper/
```

## Architecture Overview

### Three-Application System

The repository contains three separate applications that share a single PostgreSQL database:

1. **cmd/mindwell-server** - Main REST API server
   - All user-facing API endpoints
   - OAuth2 authentication and authorization
   - Real-time notifications via Centrifugo
   - Email and Telegram notifications
   - Scheduled tasks (karma recalculation, invite distribution)

2. **cmd/mindwell-images-server** - Dedicated image processing server
   - Image upload and processing
   - Async image resizing (thumbnail, small, medium, large)
   - Uses govips/v2 (libvips) for high-performance image operations
   - Separated to prevent blocking main API with I/O operations

3. **cmd/mindwell-helper** - CLI utility for batch operations
   - Administrative tasks
   - Email campaigns (reminders, surveys)
   - User activity log imports
   - Official webapp generation

### API Architecture (go-swagger based)

The entire REST API layer is generated from `web/swagger.yaml` (4476 lines) using go-swagger:

**Request Flow:**
```
cmd/mindwell-server/main.go
  ↓
restapi/configure_mindwell.go (manually edited, auto-generated foundation)
  ↓
Creates MindwellServer + calls ConfigureAPI() for each module
  ↓
Each module registers handlers: srv.API.OperationHandler = HandlerFunc(...)
```

**Core Modules** (located in `internal/app/mindwell-server/`):
- `account/` - Registration, email/password management, verification
- `users/` - Profiles, followers, themes, user listings
- `entries/` - Blog posts, feeds (live, best, friends, watching)
- `comments/` - Comment creation, feeds, visibility
- `votes/` - Voting system for entries and comments
- `favorites/` - Favorite entries management
- `watchings/` - Subscribe to entry comments
- `relations/` - Following, blocking, friend requests
- `notifications/` - Notification delivery
- `chats/` - Direct messaging
- `design/` - User theme customization
- `tags/` - Entry tags and tag-based feeds
- `badges/` - User achievements
- `complains/` - Content moderation
- `adm/` - Administrative functions
- `oauth2/` - OAuth2 token flows

Each module exports a `ConfigureAPI(srv *utils.MindwellServer)` function that wires up handlers.

### Database Layer

**Transaction Pattern:**
All database operations use `database.AutoTx` for automatic transaction management:

```go
import "github.com/sevings/mindwell-server/lib/database"

func handler(params) middleware.Responder {
    return srv.Transact(func(tx *database.AutoTx) middleware.Responder {
        tx.Query(sqlStatement, args...)
        tx.Scan(&result)
        if tx.Error() != nil {
            return operations.NewErrorDefault(500)
        }
        // Auto-commits on success or rolls back on error
        return operations.NewSuccess(result)
    })
}
```

**SQL Query Builder:**
- Uses `github.com/leporo/sqlf` for PostgreSQL type-safe query construction
- Example: `sqlf.Select(...).From(...).Where(...).Join(...)`
- Supports complex CTEs and PostgreSQL-specific features

**Key Database Features:**
- PostgreSQL with RUM extension for optimized full-text search
- Materialized views for complex aggregations (rankings, feeds)
- Stored procedures for business logic (`recalc_karma()`, `give_invites()`)
- LISTEN/NOTIFY for real-time pub/sub
- ENUM types for privacy levels, notification types, authority roles

### Real-Time Communication

**Pub/Sub Architecture:**
```
PostgreSQL LISTEN/NOTIFY
  ↓
lib/pubsub/ (event dispatcher)
  ↓
lib/notifications/ (Centrifugo, email, Telegram)
  ↓
WebSocket clients
```

Channels include: `moved_entries`, `user_badges`, `new_notification`, `chats`, etc.

### Multi-Channel Notifications

The `CompositeNotifier` pattern (in `lib/notifications/`) routes notifications through three channels:
- **Email** - SMTP via Hermes HTML templates
- **Telegram** - Bot API integration
- **Real-time** - Centrifugo WebSocket

Users can enable/disable each channel per notification type.

### Privacy and Access Control

Privacy rules are enforced at query time using complex SQL WHERE clauses:

```go
// From lib/userutil/ - AddEntryOpenQuery, AddCommentOpenQuery
// Builds conditional access based on:
// - Entry privacy level (all, registered, followers, invited, some, author)
// - User's relationship to author (follower, friend, blocked)
// - Entry visibility settings
```

This ensures users only see content they're authorized to access.

## Key Utilities and Patterns

### Core Server Object (`utils/server.go`)

The `MindwellServer` struct is the central dependency container:
- Database connection pool (`DB`)
- Configuration loader (`Config`)
- Pub/Sub listener (`PS`)
- Composite notifier (email, Telegram, Centrifugo)
- Structured loggers (API, system, email, telegram, request)
- Image URL cache (`Imgs`)
- Token hash generator (`TokenHash`)
- Rate limiters

### Handler Factory Pattern

Handlers are created using closures that capture the server instance:

```go
func newHandler(srv *utils.MindwellServer) func(params) middleware.Responder {
    return func(params operation.Params) middleware.Responder {
        // Handler implementation with access to srv
    }
}
```

This provides clean dependency injection.

### Authentication (`lib/auth/`)

OAuth2 implementation with three flows:
- **Password** - Username/password authentication
- **Authorization Code** - Standard OAuth2 flow for third-party apps
- **App** - Application-level authentication

Token management:
- Access tokens (JWT with scope-based authorization)
- Refresh tokens (database-backed)
- User loading from token: `userID := params.HTTPRequest.Context().Value("UserID").(int64)`

**Import pattern:**
```go
import libauth "github.com/sevings/mindwell-server/lib/auth"

user := libauth.LoadUser(tx, userID)
```

### Configuration (`configs/*.toml`)

Uses `zpatrick/go-config` with TOML files:
- `server.toml` - Database, email, Centrifugo, Telegram, salts, URLs
- `images.toml` - Database, image storage folder, base URL

Access via: `srv.ConfigString("section.key")`, `srv.ConfigInt()`, `srv.ConfigBool()`

### Structured Logging (`go.uber.org/zap`)

Type-specific loggers:
```go
srv.LogApi().Info("API operation", zap.String("op", "getUser"))
srv.LogSystem().Error("System error", zap.Error(err))
srv.LogEmail().Info("Email sent", zap.String("to", email))
srv.LogTelegram().Info("Telegram message", zap.Int64("chat", chatID))
srv.LogRequest().Info("Request", zap.String("method", method))
```

### Internationalization (`i18n/active.ru.toml`)

Russian localization using `nicksnyder/go-i18n/v2`:
- Error messages: `srv.NewError("error.key")`
- Localized responses throughout API
- Email templates in Russian

## Development Guidelines

### Adding New API Endpoints

1. **Update Swagger spec:** Edit `web/swagger.yaml`
2. **Regenerate code:** Run `./scripts/generate.sh`
3. **Implement handler:**
   - Create/update module in `internal/app/mindwell-server/`
   - Implement handler function
   - Register in module's `ConfigureAPI()` function
4. **Update tests:** Add integration test in `test/`

### Database Migrations

Add SQL to `scripts/update.sql` for schema changes. Complex migrations may need stored procedures.

### Privacy Rules

When querying entries or comments:
- Use `userutil.AddEntryOpenQuery()` to enforce privacy
- Use `userutil.AddCommentOpenQuery()` for comment visibility
- Consider user relationships (following, blocking, friends)

**Import:**
```go
import "github.com/sevings/mindwell-server/lib/userutil"
```

### Security Notes

- All user input is validated by go-swagger before reaching handlers
- SQL injection protection via parameterized queries (sqlf)
- XSS protection via `microcosm-cc/bluemonday` HTML sanitizer
- CSRF protection via OAuth2 token validation
- Password hashing with configurable salt
- Rate limiting on registration and API endpoints

### Image Handling

Images are handled by the separate image server:
- Clients upload to `mindwell-images-server`
- Server stores originals and generates multiple sizes asynchronously
- Image metadata stored in PostgreSQL
- Image files stored on disk at configured folder path
- URLs cached in main server for 48 hours

### Testing Strategy

Integration tests in `test/` directory:
- Use `EmailSenderMock` to avoid sending real emails
- Test full user journeys (registration, login, posting, commenting)
- Database transactions are rolled back after each test
- Run with `--failfast` to stop on first failure

### Rate Limiting

Configured via tollbooth:
- Registration: 1 request per 3600 seconds per IP
- Global: 3 requests per second per IP
- Custom error responses in JSON format

## Common Patterns

### Error Handling
```go
if tx.Error() != nil {
    return operations.NewErrorDefault(500).WithPayload(
        srv.NewError("error.database"))
}
```

### Loading Authenticated User
```go
import libauth "github.com/sevings/mindwell-server/lib/auth"

userID, ok := params.HTTPRequest.Context().Value("UserID").(int64)
if !ok {
    return operations.NewErrorUnauthorized()
}

user := libauth.LoadUser(tx, userID)
```

### Sending Notifications
```go
srv.CompNtf.NotifyNewEntry(authorID, entryID, tlogID)
// Sends via email, Telegram, and/or Centrifugo based on user preferences
```

### Publishing Real-Time Events
```go
// From database triggers or application code
srv.PS.Publish("channel_name", payload)
```

### Loading Images
```go
import "github.com/sevings/mindwell-server/lib/media"

img := media.LoadImage(srv, srv.Imgs, tx, imageID) // Cached for 48 hours
```

## File Structure Reference

```
mindwell-server/
├── cmd/                         # Application entry points
│   ├── mindwell-server/         # Main API server
│   ├── mindwell-images-server/  # Image processing server
│   └── mindwell-helper/         # CLI utility
├── internal/
│   ├── app/                     # Core application logic
│   │   ├── mindwell-server/     # Main server modules (20+ packages)
│   │   └── mindwell-images/     # Image processing service
│   └── lib/                     # Shared library packages (10 packages)
│       ├── auth/                # OAuth2 authentication
│       ├── database/            # Database transaction management
│       ├── helpers/             # App creation, config loading
│       ├── media/               # Image loading and caching
│       ├── middleware/          # HTTP logging, user activity
│       ├── notifications/       # Multi-channel notifications
│       ├── pubsub/              # PostgreSQL LISTEN/NOTIFY
│       ├── textutil/            # Text manipulation utilities
│       ├── userutil/            # User operations, privacy queries
│       └── validation/          # Email validation
├── restapi/                     # Generated API layer (main server)
├── restapi_images/              # Generated API layer (image server)
├── models/                      # Data models (~40 files)
├── utils/                       # Core server orchestrator (MindwellServer)
├── helper/                      # Helper utilities for CLI
├── web/                         # Static Swagger UI + swagger.yaml
│   └── swagger.yaml             # OpenAPI 2.0 specification (4476 lines)
├── test/                        # Integration tests
├── scripts/                     # Database and build scripts
│   ├── mindwell.sql             # Full database schema
│   ├── update.sql               # Migration scripts
│   └── generate.sh              # Code generation script
├── configs/                     # Configuration files
│   ├── server.toml              # Main server config
│   └── images.toml              # Image server config
├── i18n/                        # Internationalization (Russian)
└── go.mod                       # Go module dependencies
```

## Technology Stack Summary

- **Language:** Go 1.23.0+
- **API Framework:** go-swagger (OpenAPI 2.0)
- **Database:** PostgreSQL 14+ with RUM extension
- **SQL Builder:** leporo/sqlf
- **Real-time:** Centrifugo (WebSocket pub/sub)
- **Image Processing:** govips/v2 (libvips)
- **Authentication:** OAuth2 (golang-jwt/jwt)
- **Email:** xhit/go-simple-mail with Hermes templates
- **Notifications:** Telegram Bot API
- **Logging:** uber/zap
- **i18n:** go-i18n/v2
- **HTML Sanitization:** microcosm-cc/bluemonday
- **Markdown:** gitlab.com/golang-commonmark/markdown
- **Rate Limiting:** didip/tollbooth
