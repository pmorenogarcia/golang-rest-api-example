# Claude GitHub App Setup Guide

This document explains how to use the Claude GitHub App with your Golang REST API project.

## 📋 Overview

Two workflows are configured:

1. **[claude.yml](workflows/claude.yml)** - Interactive Claude assistant via `@claude` mentions
2. **[claude-code-review.yml](workflows/claude-code-review.yml)** - Automatic PR reviews

---

## 🔐 Setup Required

### 1. Install the Claude GitHub App

Visit: https://github.com/apps/claude-code

Install it on your repository.

### 2. Add the OAuth Token Secret

1. Go to your repository settings
2. Navigate to **Secrets and variables → Actions**
3. Click **New repository secret**
4. Name: `CLAUDE_CODE_OAUTH_TOKEN`
5. Value: Your Claude OAuth token (from https://claude.ai/github)
6. Click **Add secret**

---

## 🤖 How to Use Interactive Claude (`claude.yml`)

### Trigger Claude in Issues

Comment with `@claude` in any issue:

```markdown
@claude Can you help me implement input validation for the Pokemon endpoint?
```

### Trigger Claude in Pull Requests

1. **Comment on a PR:**
   ```markdown
   @claude Review this code and suggest improvements
   ```

2. **Reply to review comments:**
   ```markdown
   @claude How should I handle this error case?
   ```

### What Claude Can Do

Claude has access to these commands:
- `make test` - Run unit tests
- `make test-integration` - Run integration tests
- `make lint` - Run golangci-lint
- `make fmt` - Format code
- `make vet` - Run go vet
- `make build` - Build the application
- `make coverage` - Generate coverage report
- `make check` - Run all checks

### Example Use Cases

**Bug Investigation:**
```markdown
@claude The health endpoint is returning 500 errors. Can you investigate?
```

**Feature Implementation:**
```markdown
@claude Add pagination support to the Pokemon list endpoint with limit and offset parameters
```

**Code Improvement:**
```markdown
@claude Can you refactor the error handling in the Pokemon service to be more idiomatic?
```

**Testing Help:**
```markdown
@claude Write integration tests for the Pokemon endpoints
```

---

## 🔍 Automatic Code Review (`claude-code-review.yml`)

### When It Runs

Automatically triggers when:
- A new PR is opened
- New commits are pushed to an existing PR
- Changes affect `.go` files, `go.mod`, `go.sum`, `Makefile`, or workflow files

### Skip Review

Add `[skip-review]` or `[WIP]` to your PR title, or mark it as draft:
```
[WIP] Add new Pokemon endpoint
```

### What Claude Reviews

Claude automatically checks:

**✅ Code Quality**
- Idiomatic Go patterns
- Error handling
- Interface usage
- Package organization

**✅ API Design**
- RESTful principles
- Request/response models
- HTTP status codes
- Swagger documentation

**✅ Security**
- Input validation
- Error message safety
- CORS configuration

**✅ Performance**
- Goroutine safety
- Context usage
- Resource management

**✅ Testing**
- Test coverage
- Table-driven tests
- Edge cases

### Review Comments

Claude will:
- Post a comprehensive review comment on your PR
- Update the same comment on subsequent pushes (sticky comments)
- Run linters and tests if needed
- Provide specific, actionable feedback with code examples

---

## 🛠️ Configuration Details

### Go Environment

Both workflows use:
- **Go version:** 1.25.4
- **Model:** Claude Sonnet 4.5
- **Max turns:** 15-20 iterations
- **Caching:** Enabled for faster builds

### Allowed Commands

For security, Claude can only run specific make targets:
- Testing: `make test`, `make test-integration`, `make test-all`
- Quality: `make lint`, `make fmt`, `make vet`, `make check`
- Build: `make build`, `make coverage`
- Go tools: `go test`, `go vet`, `golangci-lint run`

### Customization

Edit the workflows to:
- Change allowed commands in `--allowedTools`
- Modify system prompts in `--system-prompt`
- Adjust max turns with `--max-turns`
- Change trigger patterns

---

## 📝 Best Practices

### For Interactive Use

1. **Be specific** - Give Claude context about what you need
2. **Reference files** - Mention specific files or functions
3. **Iterate** - Claude can make multiple improvements in one session
4. **Review changes** - Always review Claude's suggestions before merging

### For Code Reviews

1. **Keep PRs focused** - Smaller PRs get better reviews
2. **Write good descriptions** - Help Claude understand the context
3. **Run tests locally first** - Fix obvious issues before pushing
4. **Respond to feedback** - Claude can clarify or adjust suggestions

---

## 🚨 Troubleshooting

### Workflow Not Triggering

- Verify the GitHub App is installed on your repository
- Check the `CLAUDE_CODE_OAUTH_TOKEN` secret is set
- Ensure you're using `@claude` (not @Claude or other variants)

### Permission Errors

- Verify workflows have the correct permissions in settings
- Check: **Settings → Actions → General → Workflow permissions**
- Should be: "Read and write permissions"

### Tests Failing in Workflow

- Ensure dependencies are committed to `go.mod` and `go.sum`
- Check if tests require environment variables (add them to workflow)
- Verify tests pass locally with `make test`

---

## 📚 Additional Resources

- [Claude Code Documentation](https://docs.anthropic.com/claude/docs/claude-code)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)

---

## 🎯 Quick Start Checklist

- [ ] Install Claude GitHub App
- [ ] Add `CLAUDE_CODE_OAUTH_TOKEN` secret
- [ ] Commit both workflow files
- [ ] Push to GitHub
- [ ] Test with `@claude hello` in an issue
- [ ] Open a test PR to verify automatic reviews

---

**Need help?** Mention `@claude` in any issue and ask!
