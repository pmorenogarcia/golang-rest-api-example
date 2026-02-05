# API Usage Guide

This guide provides detailed examples of how to use the Pokemon REST API.

## Base URL

```
http://localhost:8080
```

## Authentication

Currently, this API does not require authentication. All endpoints are publicly accessible.

## Table of Contents

1. [Health Check](#health-check)
2. [List Pokemon](#list-pokemon)
3. [Get Pokemon by Name](#get-pokemon-by-name)
4. [Get Pokemon by ID](#get-pokemon-by-id)
5. [Error Responses](#error-responses)
6. [Rate Limiting](#rate-limiting)

---

## Health Check

Check if the API is running and healthy.

### Request

```bash
curl -X GET http://localhost:8080/health
```

### Response

```json
{
  "status": "ok",
  "timestamp": "2026-02-04T12:30:00Z"
}
```

**Status Code**: `200 OK`

---

## List Pokemon

Get a paginated list of Pokemon.

### Request

```bash
curl -X GET "http://localhost:8080/api/v1/pokemon?limit=5&offset=0"
```

### Query Parameters

| Parameter | Type | Required | Default | Max | Description |
|-----------|------|----------|---------|-----|-------------|
| `limit` | integer | No | 20 | 100 | Number of Pokemon to return |
| `offset` | integer | No | 0 | - | Number of Pokemon to skip |

### Response

```json
{
  "count": 1302,
  "next": "https://pokeapi.co/api/v2/pokemon?offset=5&limit=5",
  "previous": null,
  "results": [
    {
      "name": "bulbasaur",
      "url": "https://pokeapi.co/api/v2/pokemon/1/"
    },
    {
      "name": "ivysaur",
      "url": "https://pokeapi.co/api/v2/pokemon/2/"
    },
    {
      "name": "venusaur",
      "url": "https://pokeapi.co/api/v2/pokemon/3/"
    },
    {
      "name": "charmander",
      "url": "https://pokeapi.co/api/v2/pokemon/4/"
    },
    {
      "name": "charmeleon",
      "url": "https://pokeapi.co/api/v2/pokemon/5/"
    }
  ]
}
```

**Status Code**: `200 OK`

### Pagination Examples

#### Get the first page (default)
```bash
curl "http://localhost:8080/api/v1/pokemon"
```

#### Get specific page
```bash
curl "http://localhost:8080/api/v1/pokemon?limit=10&offset=20"
```

#### Get maximum allowed per page
```bash
curl "http://localhost:8080/api/v1/pokemon?limit=100&offset=0"
```

---

## Get Pokemon by Name

Get detailed information about a Pokemon using its name.

### Request

```bash
curl -X GET http://localhost:8080/api/v1/pokemon/pikachu
```

### Response

```json
{
  "id": 25,
  "name": "pikachu",
  "height": 4,
  "weight": 60,
  "base_experience": 112,
  "types": [
    {
      "slot": 1,
      "type": {
        "name": "electric",
        "url": "https://pokeapi.co/api/v2/type/13/"
      }
    }
  ],
  "abilities": [
    {
      "is_hidden": false,
      "slot": 1,
      "ability": {
        "name": "static",
        "url": "https://pokeapi.co/api/v2/ability/9/"
      }
    },
    {
      "is_hidden": true,
      "slot": 3,
      "ability": {
        "name": "lightning-rod",
        "url": "https://pokeapi.co/api/v2/ability/31/"
      }
    }
  ],
  "stats": [
    {
      "base_stat": 35,
      "effort": 0,
      "stat": {
        "name": "hp",
        "url": "https://pokeapi.co/api/v2/stat/1/"
      }
    },
    {
      "base_stat": 55,
      "effort": 0,
      "stat": {
        "name": "attack",
        "url": "https://pokeapi.co/api/v2/stat/2/"
      }
    },
    {
      "base_stat": 40,
      "effort": 0,
      "stat": {
        "name": "defense",
        "url": "https://pokeapi.co/api/v2/stat/3/"
      }
    },
    {
      "base_stat": 50,
      "effort": 0,
      "stat": {
        "name": "special-attack",
        "url": "https://pokeapi.co/api/v2/stat/4/"
      }
    },
    {
      "base_stat": 50,
      "effort": 0,
      "stat": {
        "name": "special-defense",
        "url": "https://pokeapi.co/api/v2/stat/5/"
      }
    },
    {
      "base_stat": 90,
      "effort": 2,
      "stat": {
        "name": "speed",
        "url": "https://pokeapi.co/api/v2/stat/6/"
      }
    }
  ],
  "sprites": {
    "front_default": "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/25.png",
    "front_shiny": "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/shiny/25.png",
    "back_default": "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/back/25.png",
    "back_shiny": "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/back/shiny/25.png"
  }
}
```

**Status Code**: `200 OK`

### More Examples

#### Get Charizard
```bash
curl http://localhost:8080/api/v1/pokemon/charizard
```

#### Get Mewtwo
```bash
curl http://localhost:8080/api/v1/pokemon/mewtwo
```

#### Names are case-insensitive
```bash
curl http://localhost:8080/api/v1/pokemon/PIKACHU
curl http://localhost:8080/api/v1/pokemon/PiKaChU
```

---

## Get Pokemon by ID

Get detailed information about a Pokemon using its ID.

### Request

```bash
curl -X GET http://localhost:8080/api/v1/pokemon/25
```

### Response

Same format as "Get Pokemon by Name" - returns full Pokemon details.

**Status Code**: `200 OK`

### More Examples

#### Get Bulbasaur (ID: 1)
```bash
curl http://localhost:8080/api/v1/pokemon/1
```

#### Get Mew (ID: 151)
```bash
curl http://localhost:8080/api/v1/pokemon/151
```

#### Get Rayquaza (ID: 384)
```bash
curl http://localhost:8080/api/v1/pokemon/384
```

---

## Error Responses

The API returns consistent error responses across all endpoints.

### 400 Bad Request

Invalid input or query parameters.

```bash
curl "http://localhost:8080/api/v1/pokemon?limit=200"
```

```json
{
  "error": "Bad Request",
  "message": "invalid limit: must be between 1 and 100",
  "code": 400,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### 404 Not Found

Pokemon not found.

```bash
curl http://localhost:8080/api/v1/pokemon/nonexistent
```

```json
{
  "error": "Not Found",
  "message": "Pokemon not found",
  "code": 404,
  "request_id": "550e8400-e29b-41d4-a716-446655440001"
}
```

### 500 Internal Server Error

Unexpected server error.

```json
{
  "error": "Internal Server Error",
  "message": "An unexpected error occurred",
  "code": 500,
  "request_id": "550e8400-e29b-41d4-a716-446655440002"
}
```

### 502 Bad Gateway

External API (PokeAPI) is unavailable or returned an error.

```json
{
  "error": "Bad Gateway",
  "message": "Failed to fetch data from external API",
  "code": 502,
  "request_id": "550e8400-e29b-41d4-a716-446655440003"
}
```

---

## Rate Limiting

Currently, this API does not implement rate limiting. However, please be respectful of the PokeAPI service by not making excessive requests.

**Best Practices:**
- Cache responses when possible
- Implement exponential backoff for retries
- Don't make more than 100 requests per second

---

## Using with Programming Languages

### JavaScript/Node.js

```javascript
// Using fetch
async function getPokemon(name) {
  const response = await fetch(`http://localhost:8080/api/v1/pokemon/${name}`);
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  return await response.json();
}

// Usage
getPokemon('pikachu')
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```

### Python

```python
import requests

def get_pokemon(name):
    response = requests.get(f'http://localhost:8080/api/v1/pokemon/{name}')
    response.raise_for_status()
    return response.json()

# Usage
try:
    pokemon = get_pokemon('pikachu')
    print(pokemon)
except requests.exceptions.HTTPError as e:
    print(f'Error: {e}')
```

### Go

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type Pokemon struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    // ... other fields
}

func getPokemon(name string) (*Pokemon, error) {
    resp, err := http.Get(fmt.Sprintf("http://localhost:8080/api/v1/pokemon/%s", name))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
    }

    var pokemon Pokemon
    if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
        return nil, err
    }

    return &pokemon, nil
}
```

---

## Headers

### Request Headers

The API accepts the following request headers:

- `Accept: application/json` - Expected response format (default)
- `User-Agent` - Your application identifier (optional but recommended)

### Response Headers

All responses include:

- `Content-Type: application/json`
- `X-Request-ID` - Unique request identifier for tracing

---

## Tips and Best Practices

1. **Use Pagination Wisely**: Don't request all Pokemon at once. Use reasonable page sizes (20-50).

2. **Cache Responses**: Pokemon data doesn't change frequently. Cache responses to reduce load.

3. **Handle Errors Gracefully**: Always check response status codes and handle errors appropriately.

4. **Use Request IDs**: Include the `X-Request-ID` from error responses when reporting issues.

5. **Name Normalization**: Pokemon names are case-insensitive and automatically normalized to lowercase.

6. **ID vs Name**: Using IDs is slightly faster than names, but both work equally well.

---

## Interactive Documentation

For a more interactive experience, visit the Swagger UI:

```
http://localhost:8080/swagger/index.html
```

You can test all endpoints directly from your browser!

---

## Support

If you encounter any issues or have questions:

1. Check the [README](../README.md) for setup instructions
2. Review the [Architecture Documentation](ARCHITECTURE.md) for system design
3. Open an issue on GitHub with the request ID from error responses
