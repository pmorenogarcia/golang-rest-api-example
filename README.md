# Pokemon REST API

A production-ready REST API built with Go that provides Pokemon data powered by [PokeAPI](https://pokeapi.co/). This project follows clean architecture principles and includes comprehensive OpenAPI documentation.

## Features

- **Clean Architecture**: Separation of concerns with domain, service, handler, and client layers
- **Production-Ready**: Structured logging, graceful shutdown, error handling, and middleware
- **OpenAPI/Swagger Documentation**: Interactive API documentation at `/swagger/index.html`
- **PokeAPI Integration**: Fetches real Pokemon data with retry logic and error handling
- **RESTful Design**: Standard HTTP methods and status codes
- **Configuration Management**: Environment-based configuration with sensible defaults
- **Middleware**: Request logging, panic recovery, and CORS support

## Prerequisites

- Go 1.21 or higher
- Internet connection (to fetch data from PokeAPI)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/polgarcia/golang-rest-api.git
cd golang-rest-api
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file (optional, defaults will be used):
```bash
cp .env.example .env
```

4. Edit `.env` with your preferred configuration (or use defaults).

## Running the Application

### Using Go directly:
```bash
go run cmd/api/main.go
```

### Using the run script:
```bash
./scripts/run-local.sh
```

### Using Make:
```bash
make run
```

The server will start on `http://localhost:8080` by default.

## API Endpoints

### Health Check
```
GET /health
```
Check if the API is healthy and running.

### List Pokemon
```
GET /api/v1/pokemon?limit=20&offset=0
```
Get a paginated list of Pokemon.

**Query Parameters:**
- `limit` (optional): Number of Pokemon to return (default: 20, max: 100)
- `offset` (optional): Number of Pokemon to skip (default: 0)

### Get Pokemon by Name or ID
```
GET /api/v1/pokemon/{nameOrId}
```
Get detailed information about a specific Pokemon.

**Path Parameters:**
- `nameOrId`: Pokemon name (e.g., "pikachu") or ID (e.g., "25")

### Get Pokemon Count
```
GET /api/v1/pokemon/count
```
Get the total count of Pokemon available in the PokeAPI.

### Compare Two Pokemon
```
GET /api/v1/pokemon/compare?pokemon1={name1}&pokemon2={name2}
```
Compare two Pokemon based on type effectiveness to determine which one has the advantage.

**Query Parameters:**
- `pokemon1` (required): First Pokemon name (e.g., "pikachu")
- `pokemon2` (required): Second Pokemon name (e.g., "squirtle")

### Swagger UI
```
GET /swagger/index.html
```
Interactive API documentation and testing interface.

## Example Requests

### Get Pikachu
```bash
curl http://localhost:8080/api/v1/pokemon/pikachu
```

### Get Pokemon by ID
```bash
curl http://localhost:8080/api/v1/pokemon/25
```

### List first 10 Pokemon
```bash
curl "http://localhost:8080/api/v1/pokemon?limit=10&offset=0"
```

### Get Pokemon Count
```bash
curl http://localhost:8080/api/v1/pokemon/count
```

### Compare Pokemon
```bash
curl "http://localhost:8080/api/v1/pokemon/compare?pokemon1=pikachu&pokemon2=squirtle"
```

### Health Check
```bash
curl http://localhost:8080/health
```

## Configuration

Configuration is managed through environment variables. See [.env.example](.env.example) for all available options.

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | HTTP server port | 8080 |
| `SERVER_READ_TIMEOUT` | Read timeout | 10s |
| `SERVER_WRITE_TIMEOUT` | Write timeout | 10s |
| `SERVER_IDLE_TIMEOUT` | Idle timeout | 120s |
| `POKEAPI_BASE_URL` | PokeAPI base URL | https://pokeapi.co/api/v2 |
| `POKEAPI_TIMEOUT` | PokeAPI client timeout | 30s |
| `LOG_LEVEL` | Log level (debug, info, warn, error) | info |
| `LOG_FORMAT` | Log format (json, console) | json |
| `CORS_ALLOWED_ORIGINS` | CORS allowed origins | * |

## Development

### Available Make Commands

```bash
make build            # Build the application binary
make run              # Run the application
make test             # Run unit tests
make test-integration # Run integration tests
make coverage         # Generate test coverage report
make swagger          # Generate Swagger documentation
make lint             # Run linter (golangci-lint)
make clean            # Clean build artifacts
```

### Project Structure

```
golang-rest-api/
├── cmd/api/              # Application entry point
├── internal/             # Private application code
│   ├── config/          # Configuration management
│   ├── domain/          # Domain models and interfaces
│   ├── handler/         # HTTP handlers
│   ├── service/         # Business logic
│   ├── client/          # External API clients
│   ├── middleware/      # HTTP middleware
│   └── server/          # Server setup and routing
├── pkg/                  # Public utility packages
│   └── logger/          # Structured logging
├── api/                  # OpenAPI specifications
├── docs/                 # Documentation
│   ├── swagger/         # Generated Swagger docs
│   ├── API_USAGE.md     # API usage guide
│   └── ARCHITECTURE.md  # Architecture documentation
├── test/                 # Tests
│   ├── integration/     # Integration tests
│   └── testdata/        # Test fixtures
└── scripts/              # Helper scripts
```

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for detailed architecture information.

## API Documentation

- **Interactive Swagger UI**: http://localhost:8080/swagger/index.html
- **OpenAPI Spec (JSON)**: http://localhost:8080/swagger/doc.json
- **OpenAPI Spec (YAML)**: [api/openapi.yaml](api/openapi.yaml)
- **API Usage Guide**: [docs/API_USAGE.md](docs/API_USAGE.md)

## Testing

### Run all tests:
```bash
make test
```

### Run integration tests:
```bash
make test-integration
```

### Generate coverage report:
```bash
make coverage
```

## Regenerating Swagger Documentation

After modifying handler annotations:

```bash
make swagger
# or
./scripts/generate-swagger.sh
```

## Error Handling

The API uses standard HTTP status codes:

- `200 OK`: Successful request
- `400 Bad Request`: Invalid input or parameters
- `404 Not Found`: Pokemon not found
- `500 Internal Server Error`: Server error
- `502 Bad Gateway`: External API (PokeAPI) error

Error responses follow a consistent format:
```json
{
  "error": "Not Found",
  "message": "Pokemon not found",
  "code": 404,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [PokeAPI](https://pokeapi.co/) for providing the Pokemon data
- Built with Go and following clean architecture principles

## Support

For issues, questions, or contributions, please open an issue on GitHub.

## Next Steps

This API is ready for GitHub Actions integration:
1. Automated code reviews on pull requests
2. Automatic documentation updates
3. CI/CD pipeline for testing and deployment

---

Made with Go and 💙 for Pokemon
