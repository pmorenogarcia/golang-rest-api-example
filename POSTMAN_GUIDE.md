# Postman Testing Guide

This guide will walk you through running the Pokemon REST API and testing it with Postman.

## Part 1: Starting the API Server

### Option 1: Using Make (Recommended)

```bash
cd /Users/pol.garcia/Desktop/Projects/golang-rest-api
make run
```

### Option 2: Using the Run Script

```bash
./scripts/run-local.sh
```

### Option 3: Using Go Directly

```bash
go run cmd/api/main.go
```

### What You Should See

When the server starts successfully, you'll see logs like this:

```json
{"level":"info","ts":1706964000.123,"caller":"api/main.go:48","msg":"Application starting...","log_level":"info","log_format":"json"}
{"level":"info","ts":1706964000.124,"caller":"api/main.go:59","msg":"PokeAPI client initialized","base_url":"https://pokeapi.co/api/v2","timeout":"30s"}
{"level":"info","ts":1706964000.124,"caller":"api/main.go:65","msg":"Pokemon service initialized"}
{"level":"info","ts":1706964000.124,"caller":"api/main.go:69","msg":"Handlers initialized"}
{"level":"info","ts":1706964000.124,"caller":"api/main.go:74","msg":"Server configured","port":"8080","read_timeout":"10s","write_timeout":"10s","idle_timeout":"120s"}
{"level":"info","ts":1706964000.125,"caller":"server/server.go:37","msg":"Starting HTTP server","addr":":8080"}
{"level":"info","ts":1706964000.125,"caller":"server/server.go:46","msg":"Server started successfully","addr":":8080"}
```

The server is now running at `http://localhost:8080` 🚀

---

## Part 2: Importing the Postman Collection

### Step 1: Open Postman

If you don't have Postman installed:
- Download from: https://www.postman.com/downloads/
- Or use the web version: https://web.postman.com/

### Step 2: Import the Collection

**Method A: Direct Import (Easiest)**

1. Open Postman
2. Click **"Import"** button (top-left corner)
3. Click **"Upload Files"**
4. Navigate to your project folder:
   ```
   /Users/pol.garcia/Desktop/Projects/golang-rest-api/
   ```
5. Select **`Pokemon_API.postman_collection.json`**
6. Click **"Import"**

**Method B: Drag and Drop**

1. Open Postman
2. Drag the `Pokemon_API.postman_collection.json` file directly into the Postman window
3. The collection will be imported automatically

### Step 3: Verify Import

You should now see a collection called **"Pokemon REST API"** in your left sidebar with:
- 📁 Health Check (1 request)
- 📁 Pokemon (12 requests)

### Step 4: Import Environment File (Recommended)

For easier configuration management, import the environment file:

**Method A: Import Local Environment**

1. Click the **"Import"** button again
2. Select **`Pokemon_API_Local.postman_environment.json`**
3. Click **"Import"**
4. In the top-right corner, select **"Pokemon API - Local"** from the environment dropdown

**Method B: Import Production Environment (Optional)**

1. Click **"Import"**
2. Select **`Pokemon_API_Production.postman_environment.json`**
3. Update the `base_url` value with your production URL
4. Select the environment from the dropdown

**Environment Variables Included:**
- `base_url` - Base API URL (http://localhost:8080)
- `api_version` - API version (v1)
- `health_endpoint` - Full health check URL
- `pokemon_endpoint` - Full Pokemon endpoint URL
- `swagger_url` - Swagger UI URL

**Why Use Environments?**
- ✅ Easy switching between local/production
- ✅ Centralized URL management
- ✅ Quick port changes
- ✅ Share configurations with team

---

## Part 3: Testing the API with Postman

### Quick Test: Health Check

1. In Postman's left sidebar, expand **"Pokemon REST API"** → **"Health Check"**
2. Click on **"Health Check"**
3. Click the blue **"Send"** button
4. You should see a `200 OK` response:

```json
{
    "status": "ok",
    "timestamp": "2026-02-04T12:30:00Z"
}
```

✅ **Your API is working!**

---

### Testing Individual Endpoints

#### 1. Get Pikachu

- **Request**: `Pokemon` → `Get Pokemon by Name`
- **Click**: Send
- **Expected**: `200 OK` with full Pikachu data (ID, types, abilities, stats, sprites)

#### 2. Get Pokemon by ID

- **Request**: `Pokemon` → `Get Bulbasaur (ID: 1)`
- **Click**: Send
- **Expected**: `200 OK` with Bulbasaur data

#### 3. List Pokemon

- **Request**: `Pokemon` → `List Pokemon (Default)`
- **Click**: Send
- **Expected**: `200 OK` with list of 20 Pokemon

#### 4. Test Pagination

- **Request**: `Pokemon` → `List Pokemon (First 5)`
- **Click**: Send
- **Expected**: `200 OK` with exactly 5 Pokemon
- Try changing the query parameters:
  - `limit=10&offset=0` → First 10 Pokemon
  - `limit=10&offset=10` → Next 10 Pokemon (11-20)

#### 5. Test Error Handling

- **Request**: `Pokemon` → `Pokemon Not Found (404)`
- **Click**: Send
- **Expected**: `404 Not Found` with error message

```json
{
    "error": "Not Found",
    "message": "Pokemon not found",
    "code": 404,
    "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

- **Request**: `Pokemon` → `Invalid Limit (400)`
- **Click**: Send
- **Expected**: `400 Bad Request` with validation error

---

## Part 4: Understanding the Collection Structure

### Variables

The collection uses a variable for the base URL:
- **Variable**: `{{base_url}}`
- **Value**: `http://localhost:8080`

**To change the port or host:**
1. Right-click on "Pokemon REST API" collection
2. Select "Edit"
3. Go to "Variables" tab
4. Change `base_url` value
5. Save

### Folders

**Health Check**:
- Quick endpoint to verify API is running

**Pokemon**:
- Get by name/ID
- List with pagination
- Error scenarios (404, 400)

### Example Responses

Each request includes example responses showing:
- Success cases (200 OK)
- Error cases (404, 400)
- Response structure

---

## Part 5: Advanced Testing

### Using Postman's Test Scripts

Add automatic tests to verify responses:

1. Click on any request
2. Go to **"Tests"** tab
3. Add this script:

```javascript
// Test status code
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

// Test response time
pm.test("Response time is less than 2000ms", function () {
    pm.expect(pm.response.responseTime).to.be.below(2000);
});

// Test response structure (for Pokemon endpoint)
pm.test("Response has required fields", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('id');
    pm.expect(jsonData).to.have.property('name');
    pm.expect(jsonData).to.have.property('types');
});

// Test X-Request-ID header
pm.test("X-Request-ID header is present", function () {
    pm.response.to.have.header("X-Request-ID");
});
```

4. Click **"Send"**
5. Check the **"Test Results"** tab to see passes/failures

### Running All Requests

**Collection Runner:**
1. Right-click on "Pokemon REST API" collection
2. Select **"Run collection"**
3. Select which requests to run
4. Click **"Run Pokemon REST API"**
5. See all results in a single view

---

## Part 6: Exploring Responses

### Viewing Response Data

After sending a request, you can view:

**Body Tab**:
- Pretty (formatted JSON)
- Raw (unformatted)
- Preview (rendered HTML if applicable)

**Headers Tab**:
- See response headers like:
  - `Content-Type: application/json`
  - `X-Request-ID: <uuid>`

**Cookies Tab**:
- View any cookies set by the API

### Inspecting Request Details

**Request Tab** shows:
- HTTP method (GET)
- Full URL with query parameters
- Request headers
- Request body (if applicable)

---

## Part 7: Common Testing Scenarios

### Scenario 1: Search for Your Favorite Pokemon

1. Click **"Get Pokemon by Name"**
2. Change the URL parameter from `pikachu` to any Pokemon name:
   - `charizard`
   - `mewtwo`
   - `eevee`
   - `gyarados`
   - `dragonite`
3. Click **Send**

### Scenario 2: Get All Generation 1 Pokemon

1. Click **"List Pokemon (Custom Limit)"**
2. Set query params:
   - `limit=151` (Gen 1 has 151 Pokemon)
   - `offset=0`
3. Click **Send**

**Note**: The API enforces a max limit of 100, so you'll get a 400 error. This is by design!

**Instead, do**:
- First batch: `limit=100&offset=0` (Pokemon 1-100)
- Second batch: `limit=51&offset=100` (Pokemon 101-151)

### Scenario 3: Test Edge Cases

Try these requests to see how the API handles errors:

1. **Empty name**: Change URL to `/api/v1/pokemon/` (missing name)
2. **Invalid characters**: `/api/v1/pokemon/abc123!!!`
3. **Very large ID**: `/api/v1/pokemon/999999`
4. **Negative limit**: `?limit=-5`
5. **String as limit**: `?limit=abc`

---

## Part 8: Monitoring & Debugging

### Check API Logs

While testing in Postman, watch the terminal where your API is running. You'll see:

```json
{"level":"info","msg":"Incoming request","request_id":"550e8400-...","method":"GET","path":"/api/v1/pokemon/pikachu"}
{"level":"info","msg":"Request completed","request_id":"550e8400-...","status_code":200,"duration":"450ms"}
```

### Use Request IDs for Debugging

Every response includes an `X-Request-ID` header:
1. In Postman, check the **Headers** tab of the response
2. Find `X-Request-ID`
3. Use this ID to trace the request in server logs

### Performance Testing

Watch the **Time** shown in Postman (bottom-right of response):
- Health check: ~5-10ms (no external calls)
- Get Pokemon: ~300-500ms (includes PokeAPI call)
- List Pokemon: ~100-200ms (lighter PokeAPI response)

---

## Part 9: Swagger UI Alternative

If you prefer a web-based interface, use Swagger UI:

1. With the API running, open your browser
2. Navigate to: http://localhost:8080/swagger/index.html
3. You'll see interactive API documentation
4. Click **"Try it out"** on any endpoint
5. Fill in parameters and click **"Execute"**

**Swagger UI Advantages:**
- No installation needed
- Built-in to the API
- Always up-to-date with code
- Shows request/response schemas

**Postman Advantages:**
- Richer testing capabilities
- Collection organization
- Environment variables
- Automated testing scripts

---

## Part 10: Troubleshooting

### API Not Starting?

**Error**: `bind: address already in use`
- **Solution**: Port 8080 is already in use
- **Fix**: Change `SERVER_PORT` in `.env` file or stop the other process

**Error**: `Failed to load configuration`
- **Solution**: Configuration file issue
- **Fix**: Ensure `.env.example` exists or run from project root

### Postman Connection Failed?

**Error**: `Could not get any response`
- **Check**: Is the API running? Look at your terminal
- **Check**: Is the base URL correct? Should be `http://localhost:8080`
- **Check**: Any firewall blocking the connection?

### Slow Response Times?

- **Normal**: First request is slower (DNS lookup, connection establishment)
- **Normal**: PokeAPI might be slow (external dependency)
- **Issue**: If consistently > 5 seconds, check internet connection

---

## Quick Reference Card

### API Endpoints

```
GET  /health                        → Health check
GET  /api/v1/pokemon                → List Pokemon (paginated)
GET  /api/v1/pokemon/{nameOrId}     → Get specific Pokemon
GET  /swagger/index.html            → API documentation
```

### Query Parameters

```
limit   → Number of Pokemon (1-100, default: 20)
offset  → Skip N Pokemon (default: 0)
```

### Success Status Codes

- `200 OK` → Request successful

### Error Status Codes

- `400 Bad Request` → Invalid input/parameters
- `404 Not Found` → Pokemon doesn't exist
- `500 Internal Server Error` → Server error
- `502 Bad Gateway` → PokeAPI unavailable

---

## Next Steps

1. ✅ Import the collection
2. ✅ Start the API
3. ✅ Run health check
4. ✅ Test Pokemon endpoints
5. ✅ Explore Swagger UI
6. 🎯 Build your own requests
7. 🎯 Add automated tests
8. 🎯 Create your own Pokemon app!

---

## Need Help?

- **API Logs**: Check terminal where API is running
- **Documentation**: See `README.md` and `docs/API_USAGE.md`
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **GitHub Issues**: Report bugs or ask questions

---

Happy Testing! 🚀⚡️🔥