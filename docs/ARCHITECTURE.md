# Architecture Documentation

## Overview

This Pokemon REST API follows **Clean Architecture** (also known as Hexagonal Architecture or Ports and Adapters) principles, ensuring a maintainable, testable, and scalable codebase.

## Table of Contents

1. [Architecture Principles](#architecture-principles)
2. [Layer Responsibilities](#layer-responsibilities)
3. [Dependency Flow](#dependency-flow)
4. [Folder Structure](#folder-structure)
5. [Design Patterns](#design-patterns)
6. [Data Flow](#data-flow)
7. [Error Handling Strategy](#error-handling-strategy)
8. [Technology Choices](#technology-choices)

---

## Architecture Principles

### Clean Architecture

The application is organized into concentric layers with dependencies pointing inward:

```
┌─────────────────────────────────────────────┐
│         External Interfaces                 │
│  (HTTP, CLI, External APIs)                 │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│       Presentation Layer                    │
│  (Handlers, Middleware, HTTP Routing)       │
│  internal/handler/*, internal/middleware/*  │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│       Application Layer                     │
│  (Business Logic, Use Cases)                │
│  internal/service/*                         │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│         Domain Layer (Core)                 │
│  (Entities, Interfaces, Business Rules)     │
│  internal/domain/*                          │
└─────────────────────────────────────────────┘
                   ▲
                   │
┌──────────────────┴──────────────────────────┐
│       Infrastructure Layer                  │
│  (External Clients, Database, Cache)        │
│  internal/client/*                          │
└─────────────────────────────────────────────┘
```

### Key Principles

1. **Dependency Inversion**: High-level modules don't depend on low-level modules. Both depend on abstractions.
2. **Separation of Concerns**: Each layer has a single, well-defined responsibility.
3. **Testability**: All layers can be tested independently with mocks/stubs.
4. **Independence**: Business logic is independent of frameworks, UI, and external agencies.

---

## Layer Responsibilities

### 1. Domain Layer (`internal/domain/`)

**Purpose**: Core business logic and rules.

**Components**:
- **Entities**: `Pokemon`, `PokemonList`, `PokemonType`, etc.
- **Interfaces**: `PokemonService`, `PokemonClient`
- **Domain Errors**: `ErrPokemonNotFound`, `ErrInvalidInput`

**Characteristics**:
- No external dependencies
- Pure business logic
- Defines contracts through interfaces
- Contains domain models

**Example**:
```go
// Domain interface - other layers depend on this
type PokemonService interface {
    GetByName(ctx context.Context, nameOrID string) (*Pokemon, error)
    List(ctx context.Context, limit, offset int) (*PokemonList, error)
}
```

### 2. Application Layer (`internal/service/`)

**Purpose**: Orchestrates business operations and enforces business rules.

**Components**:
- **Service Implementations**: `PokemonService`
- **Business Logic**: Validation, transformation, caching

**Characteristics**:
- Implements domain interfaces
- Contains use case logic
- Coordinates between domain and infrastructure
- Enforces business rules (e.g., max limit = 100)

**Example**:
```go
func (s *PokemonService) GetByName(ctx context.Context, nameOrID string) (*Pokemon, error) {
    // Business logic: normalize input
    nameOrID = strings.ToLower(strings.TrimSpace(nameOrID))

    // Delegate to client
    return s.client.FetchPokemon(ctx, nameOrID)
}
```

### 3. Infrastructure Layer (`internal/client/`)

**Purpose**: Interacts with external systems and services.

**Components**:
- **External Clients**: `PokeAPIClient`
- **HTTP clients, Database connections, Cache clients**

**Characteristics**:
- Implements domain interfaces
- Handles external communication
- Implements retry logic and error handling
- Transforms external data to domain models

**Example**:
```go
func (c *PokeAPIClient) FetchPokemon(ctx context.Context, nameOrID string) (*domain.Pokemon, error) {
    // HTTP request with retry logic
    // Transform external response to domain.Pokemon
}
```

### 4. Presentation Layer (`internal/handler/`, `internal/middleware/`)

**Purpose**: Handles HTTP requests and responses.

**Components**:
- **Handlers**: `GetPokemonByName`, `ListPokemon`, `HealthCheck`
- **Middleware**: Logging, recovery, CORS
- **Response helpers**: JSON serialization, error formatting

**Characteristics**:
- Converts HTTP requests to service calls
- Formats responses (JSON)
- Handles HTTP-specific concerns
- Maps domain errors to HTTP status codes

**Example**:
```go
func (h *Handler) GetPokemonByName(w http.ResponseWriter, r *http.Request) {
    nameOrID := chi.URLParam(r, "nameOrId")

    pokemon, err := h.pokemonService.GetByName(r.Context(), nameOrID)
    if err != nil {
        h.handlePokemonError(w, err) // Maps domain error to HTTP status
        return
    }

    WriteJSON(w, http.StatusOK, pokemon, h.logger)
}
```

### 5. Server Layer (`internal/server/`)

**Purpose**: HTTP server configuration and routing.

**Components**:
- **Server**: HTTP server setup, graceful shutdown
- **Routes**: Route registration, middleware chain

**Characteristics**:
- Configures HTTP server
- Registers routes and middleware
- Implements graceful shutdown

---

## Dependency Flow

```
main.go
  │
  ├──> Config (loads environment)
  │
  ├──> Logger (structured logging)
  │
  ├──> PokeAPIClient (infrastructure)
  │       │
  │       └──> implements domain.PokemonClient
  │
  ├──> PokemonService (application)
  │       │
  │       ├──> implements domain.PokemonService
  │       └──> depends on domain.PokemonClient interface
  │
  ├──> Handler (presentation)
  │       │
  │       └──> depends on domain.PokemonService interface
  │
  └──> Server
          │
          └──> uses Handler and Middleware
```

**Key Points**:
- Dependencies point inward (toward domain)
- All layers depend on domain interfaces, not implementations
- main.go is the composition root (Dependency Injection)

---

## Folder Structure

```
golang-rest-api/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point, DI container
│
├── internal/                     # Private application code
│   ├── config/                   # Configuration management
│   │   └── config.go            # Viper-based config loading
│   │
│   ├── domain/                   # Domain layer (CORE)
│   │   ├── pokemon.go           # Entities and interfaces
│   │   └── errors.go            # Domain errors
│   │
│   ├── service/                  # Application layer
│   │   └── pokemon.go           # Business logic implementation
│   │
│   ├── client/                   # Infrastructure layer
│   │   └── pokeapi.go           # External API client
│   │
│   ├── handler/                  # Presentation layer
│   │   ├── handler.go           # Base handler
│   │   ├── pokemon.go           # Pokemon endpoints
│   │   ├── health.go            # Health check
│   │   └── response.go          # Response helpers
│   │
│   ├── middleware/               # HTTP middleware
│   │   ├── logger.go            # Request logging
│   │   ├── recovery.go          # Panic recovery
│   │   └── cors.go              # CORS handling
│   │
│   └── server/                   # Server configuration
│       ├── server.go            # HTTP server setup
│       └── routes.go            # Route definitions
│
├── pkg/                          # Public packages (can be imported)
│   └── logger/
│       └── logger.go            # Zap logger wrapper
│
├── api/                          # API contracts
│   └── openapi.yaml             # OpenAPI 3.0 specification
│
├── docs/                         # Documentation
│   ├── swagger/                 # Generated Swagger docs
│   ├── API_USAGE.md            # API usage guide
│   └── ARCHITECTURE.md         # This file
│
├── test/                         # Tests
│   ├── integration/             # Integration tests
│   └── testdata/                # Test fixtures
│
├── scripts/                      # Helper scripts
│   ├── generate-swagger.sh     # Swagger generation
│   └── run-local.sh            # Local development
│
├── .env.example                 # Example environment variables
├── Makefile                     # Build automation
├── go.mod                       # Go module definition
└── README.md                    # Project README
```

---

## Design Patterns

### 1. Dependency Injection (DI)

**Pattern**: Constructor-based injection

```go
// Service depends on interface, not concrete implementation
func NewPokemonService(client domain.PokemonClient, log *logger.Logger) *PokemonService {
    return &PokemonService{
        client: client,
        logger: log,
    }
}
```

**Benefits**:
- Easy testing with mocks
- Loose coupling
- Clear dependencies

### 2. Repository Pattern

**Pattern**: Abstract data access through interfaces

```go
// Domain defines the interface
type PokemonClient interface {
    FetchPokemon(ctx context.Context, nameOrID string) (*Pokemon, error)
}

// Infrastructure implements it
type PokeAPIClient struct { ... }
func (c *PokeAPIClient) FetchPokemon(...) { ... }
```

**Benefits**:
- Swap implementations easily
- Test with mock repositories
- Hide data source details

### 3. Middleware Chain Pattern

**Pattern**: Composable request processing

```go
r.Use(middleware.Recovery(log))   // Outermost
r.Use(middleware.Logger(log))     // Middle
r.Use(middleware.CORS(origins))   // Innermost
```

**Benefits**:
- Cross-cutting concerns
- Reusable components
- Composable behavior

### 4. Error Wrapping

**Pattern**: Sentinel errors with context

```go
var ErrPokemonNotFound = errors.New("pokemon not found")

// Wrap errors for context
return nil, fmt.Errorf("%w: %v", domain.ErrExternalAPI, err)
```

**Benefits**:
- Type-safe error checking with `errors.Is()`
- Error context preservation
- Consistent error handling

---

## Data Flow

### Request Flow (Get Pokemon)

```
1. HTTP Request
   └─> GET /api/v1/pokemon/pikachu

2. Middleware Chain
   └─> Recovery → Logger → CORS

3. Handler Layer (internal/handler/pokemon.go)
   ├─> Parse path parameter: "pikachu"
   ├─> Call pokemonService.GetByName(ctx, "pikachu")
   └─> [Wait for result]

4. Service Layer (internal/service/pokemon.go)
   ├─> Validate input (not empty)
   ├─> Normalize: strings.ToLower("pikachu")
   ├─> Call client.FetchPokemon(ctx, "pikachu")
   └─> [Wait for result]

5. Client Layer (internal/client/pokeapi.go)
   ├─> Build URL: https://pokeapi.co/api/v2/pokemon/pikachu
   ├─> HTTP GET with retries (max 3)
   ├─> Parse JSON response
   ├─> Transform to domain.Pokemon
   └─> Return domain.Pokemon

6. Service Layer
   └─> Return domain.Pokemon to handler

7. Handler Layer
   ├─> Check for errors
   ├─> Format as JSON
   └─> Write HTTP 200 with JSON body

8. Middleware (Logger)
   └─> Log response: status=200, duration=150ms

9. HTTP Response
   └─> JSON response to client
```

### Error Flow

```
1. Client Layer: Pokemon not found
   └─> Return domain.ErrPokemonNotFound

2. Service Layer
   └─> Pass error through (no transformation)

3. Handler Layer
   ├─> Check error type: errors.Is(err, domain.ErrPokemonNotFound)
   ├─> Map to HTTP 404
   └─> Format error response: {"error": "Not Found", ...}

4. HTTP Response
   └─> 404 Not Found with error JSON
```

---

## Error Handling Strategy

### Domain Errors

```go
// Sentinel errors in domain layer
var (
    ErrPokemonNotFound = errors.New("pokemon not found")
    ErrInvalidInput    = errors.New("invalid input")
    ErrExternalAPI     = errors.New("external API error")
)
```

### Error Mapping

| Domain Error | HTTP Status | Description |
|--------------|-------------|-------------|
| `ErrPokemonNotFound` | 404 Not Found | Pokemon doesn't exist |
| `ErrInvalidInput` | 400 Bad Request | Validation failed |
| `ErrInvalidLimit` | 400 Bad Request | Invalid pagination limit |
| `ErrInvalidOffset` | 400 Bad Request | Invalid pagination offset |
| `ErrExternalAPI` | 502 Bad Gateway | PokeAPI unavailable |
| Other | 500 Internal Server Error | Unexpected errors |

### Error Response Format

```json
{
  "error": "Not Found",
  "message": "Pokemon not found",
  "code": 404,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

## Technology Choices

### Why These Technologies?

#### 1. **Chi Router** (github.com/go-chi/chi/v5)

**Rationale**:
- Lightweight and idiomatic Go
- Compatible with `net/http` standard library
- Excellent middleware support
- Fast and minimal dependencies

**Alternatives considered**:
- Gin: More features but heavier
- Gorilla Mux: Good but less maintained
- stdlib ServeMux: Too basic for production

#### 2. **Zap Logger** (go.uber.org/zap)

**Rationale**:
- High-performance structured logging
- Zero-allocation JSON encoding
- Production-proven (used by Uber)
- Configurable log levels

**Alternatives considered**:
- stdlib log/slog: Good but newer, less features
- Logrus: Slower performance
- Zerolog: Similar performance but less popular

#### 3. **Viper** (github.com/spf13/viper)

**Rationale**:
- Comprehensive configuration management
- Supports env vars, files, defaults
- Type-safe access
- Well-maintained

**Alternatives considered**:
- godotenv: Too simple
- envconfig: Less flexible
- Manual env parsing: Error-prone

#### 4. **Swaggo** (github.com/swaggo/swag)

**Rationale**:
- De-facto standard for Go OpenAPI
- Annotation-based (code = docs)
- Generates Swagger UI automatically
- Good community support

**Alternatives considered**:
- Manual OpenAPI spec: Hard to maintain
- go-swagger: More complex
- oapi-codegen: Different approach (spec-first)

---

## Testing Strategy

### Unit Tests

**What**: Test individual components in isolation

```go
func TestPokemonService_GetByName(t *testing.T) {
    mockClient := &MockPokemonClient{}
    service := NewPokemonService(mockClient, logger)

    pokemon, err := service.GetByName(ctx, "pikachu")
    // assertions...
}
```

**Tools**: `testify/assert`, `testify/mock`

### Integration Tests

**What**: Test full request/response cycle

```go
func TestGetPokemonEndpoint(t *testing.T) {
    server := setupTestServer()
    resp := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/api/v1/pokemon/pikachu", nil)

    server.ServeHTTP(resp, req)
    // assertions...
}
```

**Tools**: `httptest`, `testify/suite`

---

## Performance Considerations

### 1. **Context Propagation**

All operations accept `context.Context` for:
- Request cancellation
- Timeouts
- Request tracing

### 2. **HTTP Client Configuration**

```go
httpClient: &http.Client{
    Timeout: 30 * time.Second,  // Prevent hanging
}
```

### 3. **Retry Logic**

- Exponential backoff
- Max 3 retries
- Don't retry on 404

### 4. **Graceful Shutdown**

- 30-second grace period
- Finish in-flight requests
- Clean resource cleanup

---

## Security Considerations

### 1. **Input Validation**

- Validate all user inputs
- Limit pagination (max 100)
- Normalize strings (lowercase)

### 2. **Error Information Disclosure**

- Don't expose internal errors
- Use generic messages for 500 errors
- Log detailed errors server-side

### 3. **CORS Configuration**

- Configurable allowed origins
- Default: `*` (development only)
- Production: Specific origins

### 4. **Request Timeouts**

- Read timeout: 10s
- Write timeout: 10s
- Idle timeout: 120s

---

## Scalability Considerations

### Current Architecture

- Stateless application (horizontal scaling ready)
- No shared state between requests
- External API client with retry logic

### Future Enhancements

1. **Caching Layer**
   - In-memory cache (TTL: 1 hour)
   - Redis for distributed caching

2. **Database Layer**
   - Cache frequently accessed Pokemon
   - Reduce PokeAPI load

3. **Metrics & Monitoring**
   - Prometheus metrics
   - Request latency tracking
   - Error rate monitoring

4. **Rate Limiting**
   - Per-IP rate limiting
   - Protect against abuse

---

## Maintenance & Evolution

### Adding a New Endpoint

1. Add domain models/interfaces (if needed)
2. Implement service logic
3. Add handler with Swagger annotations
4. Register route in `routes.go`
5. Regenerate Swagger docs
6. Update documentation

### Changing External API

1. Create new client implementing `domain.PokemonClient`
2. Update DI in `main.go`
3. No changes needed in other layers!

### Adding a Feature

1. Start with domain layer (what data/interfaces?)
2. Implement service logic
3. Add infrastructure (external clients, DB)
4. Expose through handlers
5. Update docs and tests

---

## Conclusion

This architecture provides:

- **Maintainability**: Clear separation of concerns
- **Testability**: Mockable interfaces at every layer
- **Scalability**: Stateless, horizontally scalable
- **Flexibility**: Easy to swap implementations
- **Developer Experience**: Clean, idiomatic Go code

The architecture is production-ready and can evolve with your needs while maintaining code quality and clarity.
