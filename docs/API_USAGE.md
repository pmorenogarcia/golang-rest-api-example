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
5. [Get Pokemon Count](#get-pokemon-count)
6. [Compare Two Pokemon](#compare-two-pokemon)
7. [Error Responses](#error-responses)
8. [Rate Limiting](#rate-limiting)

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

## Get Pokemon Count

Get the total number of Pokemon available in the PokeAPI.

### Request

```bash
curl -X GET http://localhost:8080/api/v1/pokemon/count
```

### Response

```json
{
  "count": 1302
}
```

**Status Code**: `200 OK`

### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `count` | integer | Total number of Pokemon available |

### Use Cases

This endpoint is useful for:
- Displaying total Pokemon count in your UI
- Validating Pokemon IDs before making requests
- Implementing pagination (knowing total pages)
- Analytics and statistics

### Example Integration

#### JavaScript/Node.js
```javascript
async function getPokemonCount() {
  const response = await fetch('http://localhost:8080/api/v1/pokemon/count');
  const data = await response.json();
  console.log(`Total Pokemon: ${data.count}`);
  return data.count;
}
```

#### Python
```python
import requests

def get_pokemon_count():
    response = requests.get('http://localhost:8080/api/v1/pokemon/count')
    response.raise_for_status()
    data = response.json()
    print(f"Total Pokemon: {data['count']}")
    return data['count']
```

#### Go
```go
type PokemonCount struct {
    Count int `json:"count"`
}

func getPokemonCount() (int, error) {
    resp, err := http.Get("http://localhost:8080/api/v1/pokemon/count")
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    var count PokemonCount
    if err := json.NewDecoder(resp.Body).Decode(&count); err != nil {
        return 0, err
    }

    return count.Count, nil
}
```

---

## Compare Two Pokemon

Compare two Pokemon based on their type effectiveness to determine which one has the advantage.

### Request

```bash
curl -X GET "http://localhost:8080/api/v1/pokemon/compare?pokemon1=pikachu&pokemon2=squirtle"
```

### Query Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `pokemon1` | string | Yes | Name of the first Pokemon (case-insensitive) |
| `pokemon2` | string | Yes | Name of the second Pokemon (case-insensitive) |

### Response - Type Advantage

When one Pokemon has a type advantage over the other:

```json
{
  "message": "Pikachu is stronger than Squirtle because electric type beats water type",
  "type_advantage": "electric beats water",
  "winner": {
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
    "abilities": [],
    "stats": [],
    "sprites": {
      "front_default": "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/25.png",
      "front_shiny": "",
      "back_default": "",
      "back_shiny": ""
    }
  },
  "pokemon1": {
    "id": 25,
    "name": "pikachu",
    "types": [
      {
        "slot": 1,
        "type": {
          "name": "electric",
          "url": "https://pokeapi.co/api/v2/type/13/"
        }
      }
    ]
  },
  "pokemon2": {
    "id": 7,
    "name": "squirtle",
    "types": [
      {
        "slot": 1,
        "type": {
          "name": "water",
          "url": "https://pokeapi.co/api/v2/type/11/"
        }
      }
    ]
  }
}
```

**Status Code**: `200 OK`

### Response - Neutral Matchup

When neither Pokemon has a type advantage:

```json
{
  "message": "Neither Pikachu (electric) nor Charmander (fire) has a type advantage - it's a neutral matchup!",
  "winner": null,
  "pokemon1": {
    "id": 25,
    "name": "pikachu",
    "types": [
      {
        "slot": 1,
        "type": {
          "name": "electric",
          "url": "https://pokeapi.co/api/v2/type/13/"
        }
      }
    ]
  },
  "pokemon2": {
    "id": 4,
    "name": "charmander",
    "types": [
      {
        "slot": 1,
        "type": {
          "name": "fire",
          "url": "https://pokeapi.co/api/v2/type/10/"
        }
      }
    ]
  }
}
```

**Status Code**: `200 OK`

### Response - Mutual Advantage (Tie)

When both Pokemon have type advantages against each other:

```json
{
  "message": "Charizard and Golem both have type advantages against each other - it's a tie!",
  "type_advantage": "flying beats ground, but rock beats fire",
  "winner": null,
  "pokemon1": {
    "id": 6,
    "name": "charizard",
    "types": [
      {
        "slot": 1,
        "type": {
          "name": "fire",
          "url": "https://pokeapi.co/api/v2/type/10/"
        }
      },
      {
        "slot": 2,
        "type": {
          "name": "flying",
          "url": "https://pokeapi.co/api/v2/type/3/"
        }
      }
    ]
  },
  "pokemon2": {
    "id": 76,
    "name": "golem",
    "types": [
      {
        "slot": 1,
        "type": {
          "name": "rock",
          "url": "https://pokeapi.co/api/v2/type/6/"
        }
      },
      {
        "slot": 2,
        "type": {
          "name": "ground",
          "url": "https://pokeapi.co/api/v2/type/5/"
        }
      }
    ]
  }
}
```

**Status Code**: `200 OK`

### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Human-readable description of the comparison result |
| `type_advantage` | string | Description of which type beats which (optional, only when advantage exists) |
| `winner` | object/null | Full Pokemon object of the winner, or `null` if tie/neutral |
| `pokemon1` | object | Full Pokemon object for the first Pokemon |
| `pokemon2` | object | Full Pokemon object for the second Pokemon |

### Comparison Examples

#### Electric vs Water (Type Advantage)
```bash
curl "http://localhost:8080/api/v1/pokemon/compare?pokemon1=pikachu&pokemon2=squirtle"
# Result: Pikachu wins (electric beats water)
```

#### Fire vs Water (Type Advantage)
```bash
curl "http://localhost:8080/api/v1/pokemon/compare?pokemon1=charmander&pokemon2=squirtle"
# Result: Squirtle wins (water beats fire)
```

#### Grass vs Fire (Type Advantage)
```bash
curl "http://localhost:8080/api/v1/pokemon/compare?pokemon1=bulbasaur&pokemon2=charmander"
# Result: Charmander wins (fire beats grass)
```

#### Electric vs Fire (Neutral)
```bash
curl "http://localhost:8080/api/v1/pokemon/compare?pokemon1=pikachu&pokemon2=charmander"
# Result: Neutral matchup (no type advantage)
```

#### Names are case-insensitive
```bash
curl "http://localhost:8080/api/v1/pokemon/compare?pokemon1=PIKACHU&pokemon2=SQUIRTLE"
curl "http://localhost:8080/api/v1/pokemon/compare?pokemon1=PiKaChU&pokemon2=sQuIrTlE"
```

### Type Effectiveness Rules

The comparison is based on Pokemon type effectiveness:

**Super Effective Matchups:**
- Water beats Fire, Ground, Rock
- Fire beats Grass, Ice, Bug, Steel
- Grass beats Water, Ground, Rock
- Electric beats Water, Flying
- Ground beats Electric, Fire, Poison, Rock, Steel
- Rock beats Fire, Ice, Flying, Bug
- Ice beats Grass, Ground, Flying, Dragon
- Fighting beats Normal, Ice, Rock, Dark, Steel
- Poison beats Grass, Fairy
- Flying beats Grass, Fighting, Bug
- Psychic beats Fighting, Poison
- Bug beats Grass, Psychic, Dark
- Ghost beats Psychic, Ghost
- Dragon beats Dragon
- Dark beats Psychic, Ghost
- Steel beats Ice, Rock, Fairy
- Fairy beats Fighting, Dragon, Dark

### Error Cases

#### Missing Parameters
```bash
curl "http://localhost:8080/api/v1/pokemon/compare?pokemon1=pikachu"
```

```json
{
  "error": "Bad Request",
  "message": "both Pokemon names must be provided",
  "code": 400,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

#### Same Pokemon
```bash
curl "http://localhost:8080/api/v1/pokemon/compare?pokemon1=pikachu&pokemon2=pikachu"
```

```json
{
  "error": "Bad Request",
  "message": "cannot compare a Pokemon with itself",
  "code": 400,
  "request_id": "550e8400-e29b-41d4-a716-446655440001"
}
```

#### Pokemon Not Found
```bash
curl "http://localhost:8080/api/v1/pokemon/compare?pokemon1=pikachu&pokemon2=nonexistent"
```

```json
{
  "error": "Not Found",
  "message": "Pokemon not found",
  "code": 404,
  "request_id": "550e8400-e29b-41d4-a716-446655440002"
}
```

### Use Cases

This endpoint is useful for:
- Building Pokemon battle simulators
- Creating type matchup calculators
- Educational tools for learning type effectiveness
- Game strategy planning
- Pokemon team composition analysis

### Example Integration

#### JavaScript/Node.js
```javascript
async function comparePokemon(pokemon1, pokemon2) {
  const url = `http://localhost:8080/api/v1/pokemon/compare?pokemon1=${pokemon1}&pokemon2=${pokemon2}`;
  const response = await fetch(url);

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  const data = await response.json();
  console.log(data.message);

  if (data.winner) {
    console.log(`Winner: ${data.winner.name}`);
  } else {
    console.log('No clear winner!');
  }

  return data;
}

// Usage
comparePokemon('pikachu', 'squirtle')
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```

#### Python
```python
import requests

def compare_pokemon(pokemon1, pokemon2):
    url = f'http://localhost:8080/api/v1/pokemon/compare'
    params = {'pokemon1': pokemon1, 'pokemon2': pokemon2}

    response = requests.get(url, params=params)
    response.raise_for_status()

    data = response.json()
    print(data['message'])

    if data.get('winner'):
        print(f"Winner: {data['winner']['name']}")
    else:
        print('No clear winner!')

    return data

# Usage
try:
    result = compare_pokemon('pikachu', 'squirtle')
    print(result)
except requests.exceptions.HTTPError as e:
    print(f'Error: {e}')
```

#### Go
```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
)

type ComparisonResult struct {
    Message       string   `json:"message"`
    TypeAdvantage string   `json:"type_advantage,omitempty"`
    Winner        *Pokemon `json:"winner"`
    Pokemon1      *Pokemon `json:"pokemon1"`
    Pokemon2      *Pokemon `json:"pokemon2"`
}

func comparePokemon(pokemon1, pokemon2 string) (*ComparisonResult, error) {
    baseURL := "http://localhost:8080/api/v1/pokemon/compare"
    params := url.Values{}
    params.Add("pokemon1", pokemon1)
    params.Add("pokemon2", pokemon2)

    fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

    resp, err := http.Get(fullURL)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
    }

    var result ComparisonResult
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &result, nil
}

// Usage
func main() {
    result, err := comparePokemon("pikachu", "squirtle")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Println(result.Message)
    if result.Winner != nil {
        fmt.Printf("Winner: %s\n", result.Winner.Name)
    }
}
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
