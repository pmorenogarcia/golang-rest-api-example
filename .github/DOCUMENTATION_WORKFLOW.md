# 📚 Documentation Automation Workflow

This guide explains how Claude automatically helps maintain API documentation when endpoints change.

## 🔄 Workflow Overview

### Phase 1: PR Review (Automatic)
When you create a PR with endpoint changes:

1. **Claude detects endpoint changes** by analyzing:
   - `internal/server/routes.go` - new routes
   - `internal/handler/*.go` - new handlers
   - Modified or deleted endpoints

2. **Claude mentions in review**:
   ```markdown
   ## 📚 Documentation Updates Needed
   This PR adds/modifies/deletes API endpoints. Documentation should be updated.

   **Changed Endpoints:**
   - GET /api/v1/pokemon/random (new)

   **To update documentation, comment:** `@claude update docs`
   ```

### Phase 2: Documentation Update (User-Triggered)
When you're ready for Claude to update documentation:

1. **Comment on the PR with the dedicated command:**
   ```
   @claude-update-docs
   ```

   *Alternative (natural language):*
   ```
   @claude update docs
   ```

2. **Claude automatically:**
   - Reads current Postman collection from `docs/postman/`
   - Analyzes ALL endpoint changes in the PR
   - Generates updated Postman collection with:
     - New requests for new endpoints
     - Updated requests for modified endpoints
     - Removed requests for deleted endpoints
     - Proper examples, descriptions, and test cases
   - Updates `README.md` if API documentation exists there
   - Creates a commit: `docs: update API documentation for [endpoint changes]`
   - Pushes directly to your PR branch

3. **You review the changes:**
   - Check the new commit in your PR
   - Verify Postman collection is correct
   - Test the updated documentation
   - Merge when satisfied

---

## 📁 Files That Get Updated

### Automatically Updated:
- ✅ `docs/postman/Pokemon_API.postman_collection.json` - Postman collection
- ✅ `docs/postman/Pokemon_API_Local.postman_environment.json` - If environment vars needed
- ✅ `README.md` - If it contains API documentation
- ✅ Any other `*.md` files with endpoint documentation

### Never Modified:
- ❌ Source code files (`*.go`)
- ❌ Tests
- ❌ Configuration files

---

## 🎯 Example Workflow

### Scenario: Adding a new endpoint

**1. Create your PR with new endpoint:**
```go
// internal/server/routes.go
r.Get("/random", h.GetRandomPokemon)
```

**2. Claude reviews and detects:**
```markdown
## 📚 Documentation Updates Needed

**Changed Endpoints:**
- GET /api/v1/pokemon/random (new) - Returns a random Pokemon

**To update documentation, comment:** `@claude update docs`
```

**3. You trigger documentation update:**
```
@claude update docs
```

**4. Claude commits:**
```
docs: update API documentation for random Pokemon endpoint

- Add Postman request for GET /api/v1/pokemon/random
- Include success and error examples
- Update API documentation in README.md
```

**5. Check the commit and merge!**

---

## 🔧 Customization

### Change Trigger Phrase
Edit `.github/workflows/claude.yml`:
```yaml
# Change "@claude update docs" to your preferred phrase
# Look for the system prompt and update the instruction
```

### Add More Documentation Files
Claude will automatically detect and update:
- Any file in `docs/` directory
- `README.md`, `API.md`, `ENDPOINTS.md`
- Swagger/OpenAPI specs if present

### Disable Auto-Documentation
To disable this feature:
1. Remove the documentation section from the system prompt in `claude.yml`
2. Or simply don't trigger with `@claude update docs`

---

## 🚨 Troubleshooting

### Claude doesn't detect endpoint changes
**Check:**
- Did you modify `internal/server/routes.go` or handlers?
- Did Claude finish the review?
- Try asking explicitly: `@claude did I add any new endpoints?`

### Claude doesn't commit documentation
**Check:**
- Did you use the exact phrase: `@claude update docs`?
- Does the workflow have write permissions? (Check `.github/workflows/claude.yml`)
- Are there actual endpoint changes in the PR?

### Documentation commit fails
**Check:**
- GitHub Actions permissions in repository settings
- Workflow has `contents: write` permission
- Branch protection rules don't block bot commits

### Postman collection is incorrect
**Provide feedback:**
```
@claude The Postman request for /random endpoint is incorrect.
The response should include the Pokemon object, not just the ID.
Please fix and update docs again.
```

---

## 💡 Pro Tips

### 1. Review Before Triggering
Review your code changes first, then trigger docs update after you're happy with the implementation.

### 2. Batch Multiple Changes
If you're adding multiple endpoints in one PR, Claude will update documentation for all of them in a single commit.

### 3. Customize Postman Examples
After Claude commits, you can manually adjust Postman examples if needed, then commit your tweaks.

### 4. Test the Documentation
Use the updated Postman collection to test your new endpoints before merging the PR.

### 5. Ask for Specific Updates
Be specific if you want only certain docs updated:
```
@claude update only the Postman collection, skip README
```

---

## 📊 What Claude Checks

### For New Endpoints:
- ✅ HTTP method (GET, POST, PUT, DELETE, etc.)
- ✅ URL path and parameters
- ✅ Request body structure
- ✅ Response body structure
- ✅ Status codes
- ✅ Error responses
- ✅ Authentication requirements
- ✅ Example requests and responses

### For Modified Endpoints:
- ✅ Changed parameters
- ✅ Changed request/response format
- ✅ New status codes
- ✅ Updated descriptions

### For Deleted Endpoints:
- ✅ Removes from Postman collection
- ✅ Removes from documentation
- ✅ Adds deprecation note if partial removal

---

## 🎨 Postman Collection Structure

Claude maintains this structure:
```json
{
  "info": { ... },
  "item": [
    {
      "name": "Pokemon",
      "item": [
        {
          "name": "Get Pokemon by Name",
          "request": { ... },
          "response": [ ... ]
        },
        {
          "name": "Get Random Pokemon",  // ← Claude adds this
          "request": { ... },
          "response": [ ... ]
        }
      ]
    }
  ]
}
```

---

## ✅ Best Practices

1. **Always review docs commits** - Claude is good but verify the changes
2. **Trigger after code is stable** - Don't update docs while still iterating on code
3. **Use descriptive endpoint names** - Helps Claude generate better documentation
4. **Add Swagger comments** - Claude uses these for better docs generation
5. **Test with Postman** - Import the updated collection and test your endpoints

---

## 🔒 Security Notes

- Claude only commits to YOUR PR branch, never directly to main
- Documentation commits go through the same review process as code
- You can always revert documentation commits if needed
- Claude never exposes secrets or sensitive data in documentation

---

## 📞 Getting Help

If something doesn't work:
1. Check this guide's troubleshooting section
2. Ask Claude directly: `@claude help with documentation workflow`
3. Check GitHub Actions logs for the workflow run
4. Open an issue if it's a recurring problem

---

**Happy documenting! 🚀**
