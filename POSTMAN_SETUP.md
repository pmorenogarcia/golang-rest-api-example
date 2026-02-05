# Postman Quick Setup Guide

## 📦 Files to Import

You'll find these files in the project root:

### 1. Collection File (Required)
**File**: `Pokemon_API.postman_collection.json`
- Contains all 13 pre-configured API requests
- Organized into folders (Health Check, Pokemon)
- Includes example responses

### 2. Environment Files (Recommended)

**Local Environment**: `Pokemon_API_Local.postman_environment.json`
- Pre-configured for `http://localhost:8080`
- Use this for local development

**Production Environment**: `Pokemon_API_Production.postman_environment.json`
- Template for production deployment
- Update `base_url` with your actual domain

---

## 🚀 Quick Start (4 Steps)

### Step 1: Start the API
```bash
cd /Users/pol.garcia/Desktop/Projects/golang-rest-api
make run
```

**Verify it's running:**
- Open browser: http://localhost:8080/health
- Should see: `{"status":"ok","timestamp":"..."}`

---

### Step 2: Import Collection into Postman

1. Open **Postman**
2. Click **"Import"** (top-left)
3. Drag and drop **`Pokemon_API.postman_collection.json`**
4. Collection appears in left sidebar

---

### Step 3: Import Environment (Recommended)

1. Click **"Import"** again
2. Drag and drop **`Pokemon_API_Local.postman_environment.json`**
3. In top-right corner, select **"Pokemon API - Local"** from dropdown
4. ✅ Environment is now active

---

### Step 4: Test Your First Request

1. Expand **"Pokemon REST API"** → **"Health Check"**
2. Click **"Health Check"**
3. Click blue **"Send"** button
4. See `200 OK` response ✅

---

## 🎯 What's in the Collection?

### Health Check (1 request)
- ✅ Health Check

### Pokemon Endpoints (12 requests)

**Get Pokemon:**
- ✅ Get Pokemon by Name (Pikachu)
- ✅ Get Pokemon by ID (25)
- ✅ Get Charizard
- ✅ Get Mewtwo
- ✅ Get Bulbasaur (ID: 1)

**List Pokemon:**
- ✅ List Pokemon (Default - 20 results)
- ✅ List Pokemon (Custom Limit - 10 results)
- ✅ List Pokemon (Page 2 - Offset pagination)
- ✅ List Pokemon (First 5)

**Error Testing:**
- ✅ Pokemon Not Found (404 error)
- ✅ Invalid Limit (400 error)
- ✅ Negative Offset (400 error)

---

## 🔧 Environment Variables

When you import the environment, these variables are available:

| Variable | Value | Usage |
|----------|-------|-------|
| `base_url` | `http://localhost:8080` | Main API URL |
| `api_version` | `v1` | API version |
| `health_endpoint` | `{{base_url}}/health` | Health check URL |
| `pokemon_endpoint` | `{{base_url}}/api/{{api_version}}/pokemon` | Pokemon endpoint |
| `swagger_url` | `{{base_url}}/swagger/index.html` | Swagger UI |

**How to use in requests:**
- Use `{{base_url}}` in any request URL
- Change port: Edit `base_url` in environment (e.g., `http://localhost:3000`)
- Switch environments: Use dropdown in top-right corner

---

## 🧪 Quick Tests

Try these requests right away:

### 1. Health Check ✅
```
GET {{base_url}}/health
```
Expected: `200 OK`

### 2. Get Pikachu ⚡
```
GET {{base_url}}/api/v1/pokemon/pikachu
```
Expected: `200 OK` with Pikachu data

### 3. List First 5 Pokemon 📋
```
GET {{base_url}}/api/v1/pokemon?limit=5&offset=0
```
Expected: `200 OK` with 5 Pokemon

### 4. Test 404 Error ❌
```
GET {{base_url}}/api/v1/pokemon/nonexistent
```
Expected: `404 Not Found`

---

## 🌐 Changing the Base URL

### Method 1: Edit Environment (Recommended)

1. Click environment dropdown (top-right)
2. Click the eye icon 👁️
3. Click **"Edit"**
4. Change `base_url` value
5. **Save**

### Method 2: Edit Collection Variable

1. Right-click **"Pokemon REST API"** collection
2. Select **"Edit"**
3. Go to **"Variables"** tab
4. Change `base_url` initial value
5. **Save**

---

## 📊 Response Headers to Check

Every response includes useful headers:

- `Content-Type: application/json` - Response format
- `X-Request-ID: <uuid>` - Unique request identifier for tracing
- `Access-Control-Allow-Origin: *` - CORS enabled

**Use Request ID for debugging:**
1. Copy `X-Request-ID` from response headers
2. Search for it in API server logs
3. Trace the full request lifecycle

---

## 🎮 Popular Pokemon to Try

Change the URL parameter to any of these:

**Starters:**
- `/api/v1/pokemon/bulbasaur` (001)
- `/api/v1/pokemon/charmander` (004)
- `/api/v1/pokemon/squirtle` (007)

**Legendary:**
- `/api/v1/pokemon/articuno` (144)
- `/api/v1/pokemon/zapdos` (145)
- `/api/v1/pokemon/moltres` (146)
- `/api/v1/pokemon/mewtwo` (150)
- `/api/v1/pokemon/mew` (151)

**Fan Favorites:**
- `/api/v1/pokemon/pikachu` (025)
- `/api/v1/pokemon/eevee` (133)
- `/api/v1/pokemon/dragonite` (149)
- `/api/v1/pokemon/gyarados` (130)
- `/api/v1/pokemon/charizard` (006)

---

## 🔍 Testing Edge Cases

Try these to see error handling:

### Validation Errors (400)
```
GET {{base_url}}/api/v1/pokemon?limit=200    # Exceeds max
GET {{base_url}}/api/v1/pokemon?limit=-5     # Negative limit
GET {{base_url}}/api/v1/pokemon?offset=-1    # Negative offset
GET {{base_url}}/api/v1/pokemon?limit=abc    # Invalid type
```

### Not Found Errors (404)
```
GET {{base_url}}/api/v1/pokemon/notarealmon
GET {{base_url}}/api/v1/pokemon/999999
GET {{base_url}}/api/v1/pokemon/xyz123
```

---

## 🚨 Troubleshooting

### Problem: "Could not get any response"

**Solutions:**
1. ✅ Check if API is running: `ps aux | grep pokemon-api`
2. ✅ Verify URL: Should be `http://localhost:8080` (not https)
3. ✅ Check port: Is something else using port 8080?
4. ✅ Environment selected: See top-right dropdown

### Problem: "404 Not Found" for valid endpoint

**Solutions:**
1. ✅ Check base URL has no trailing slash
2. ✅ Verify path: `/api/v1/pokemon` (not `/api/pokemon`)
3. ✅ Case-sensitive: Use lowercase for Pokemon names

### Problem: "Connection refused"

**Solutions:**
1. ✅ Start the API: `make run`
2. ✅ Check terminal for errors
3. ✅ Verify port in `.env` file

### Problem: Slow responses (> 5 seconds)

**Causes:**
- First request is always slower (DNS, connection setup)
- PokeAPI (external service) might be slow
- Check your internet connection

---

## 📖 Complete Documentation

For more details, see:
- **POSTMAN_GUIDE.md** - Full testing guide (10 parts)
- **API_USAGE.md** - API documentation with examples
- **README.md** - Project overview

---

## 🎯 Alternative: Swagger UI

Don't want to use Postman? Try Swagger UI:

1. Start API: `make run`
2. Open: http://localhost:8080/swagger/index.html
3. Click **"Try it out"** on any endpoint
4. Test directly in browser

**Swagger Advantages:**
- No installation
- Always up-to-date
- Built-in to API

**Postman Advantages:**
- Advanced testing
- Collection organization
- Environment management
- Team collaboration

---

## 📝 Quick Reference Card

### Files to Import
```
✅ Pokemon_API.postman_collection.json       (Required)
✅ Pokemon_API_Local.postman_environment.json (Recommended)
```

### API Endpoints
```
GET /health                        → Health check
GET /api/v1/pokemon                → List Pokemon
GET /api/v1/pokemon/{nameOrId}     → Get specific Pokemon
GET /swagger/index.html            → API docs
```

### Status Codes
```
200 OK          → Success
400 Bad Request → Invalid input
404 Not Found   → Pokemon doesn't exist
500 Server Error → Internal error
502 Bad Gateway  → PokeAPI unavailable
```

---

## ✅ Checklist

Before testing:
- [ ] API is running (`make run`)
- [ ] Collection imported
- [ ] Environment imported
- [ ] Environment selected (top-right dropdown)
- [ ] Health check works

Ready to test:
- [ ] Get Pokemon by name
- [ ] Get Pokemon by ID
- [ ] List Pokemon with pagination
- [ ] Test error cases (404, 400)
- [ ] Check request/response times

---

## 🎉 You're All Set!

Everything is configured and ready to go. Start testing your Pokemon API!

**Next Steps:**
1. 🧪 Test all 13 endpoints
2. 🎮 Search for your favorite Pokemon
3. 📊 Run Collection Runner for automated testing
4. 🚀 Build something awesome!

**Need help?** Check the full guide: `POSTMAN_GUIDE.md`

---

Happy Testing! ⚡️🔥🌊