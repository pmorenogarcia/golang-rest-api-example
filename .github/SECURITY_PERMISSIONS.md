# 🔒 Claude Security Permissions

This document explains the security restrictions in place for Claude's repository access.

## 📋 Permission Model

### Read Access (Entire Repository)
Claude can **READ** any file in the repository:
- ✅ Source code (`*.go`, `internal/`, `cmd/`, `pkg/`)
- ✅ Tests (`test/`, `*_test.go`)
- ✅ Configuration (`.github/`, `go.mod`, `Makefile`)
- ✅ Documentation (`docs/`, `README.md`)
- ✅ All other files

### Write Access (Restricted to Documentation)
Claude can **WRITE** only to:
- ✅ `docs/` folder (entire directory and all subdirectories)
- ✅ `README.md` (root level only)
- ✅ `API.md` (root level only)
- ✅ `ENDPOINTS.md` (root level only)
- ✅ `CHANGELOG.md` (root level only)

### Forbidden Write Access
Claude **CANNOT** write to:
- ❌ Source code (`*.go` files)
- ❌ Internal packages (`internal/`, `cmd/`, `pkg/`)
- ❌ Tests (`test/`, `*_test.go`)
- ❌ Build configuration (`Makefile`, `go.mod`, `go.sum`)
- ❌ CI/CD workflows (`.github/workflows/`)
- ❌ Git configuration (`.gitignore`, `.git/`)
- ❌ Any file outside the allowed write list

---

## 🛡️ Enforcement Mechanism

### 1. Pre-commit Hook
Location: `.git/hooks/pre-commit`

Installed automatically by the workflow, this hook:
- Runs before every commit Claude attempts
- Checks all staged files
- Validates each file is in an allowed write location
- **BLOCKS** the commit if forbidden files are detected

### 2. Validation Script
Location: `.github/scripts/validate-docs-only.sh`

This bash script:
```bash
#!/bin/bash
# Validates that only docs/ folder files are being committed

# Check each staged file
# Allow: docs/**, README.md, API.md, ENDPOINTS.md, CHANGELOG.md
# Block: Everything else
```

### 3. System Prompt Restrictions
The Claude AI system prompt explicitly states:
- Write access limited to docs/ folder
- Pre-commit hook will block violations
- Instructions to only stage docs files

---

## 🎯 How It Works

### Scenario 1: Claude Tries to Commit Documentation ✅

```bash
# Claude runs:
git add docs/postman/Pokemon_API.postman_collection.json
git add README.md
git commit -m "docs: update API documentation"
```

**Result:**
```
🔍 Validating that only documentation files are being committed...
  ✅ docs/postman/Pokemon_API.postman_collection.json (docs folder)
  ✅ README.md (allowed documentation file)

✅ All staged files are valid documentation files
[main abc1234] docs: update API documentation
 2 files changed, 45 insertions(+)
```

### Scenario 2: Claude Tries to Commit Source Code ❌

```bash
# Claude runs:
git add internal/handler/pokemon.go
git add docs/postman/Pokemon_API.postman_collection.json
git commit -m "fix: update handler and docs"
```

**Result:**
```
🔍 Validating that only documentation files are being committed...
  ❌ internal/handler/pokemon.go (NOT in docs folder or allowed list)
  ✅ docs/postman/Pokemon_API.postman_collection.json (docs folder)

❌ ERROR: Claude can only commit to the docs/ folder!

Invalid files detected:
  - internal/handler/pokemon.go

Allowed locations:
  ✅ docs/**/*
  ✅ README.md, API.md, ENDPOINTS.md, CHANGELOG.md

Please remove non-documentation files from this commit.
```

**Commit BLOCKED** - does not proceed.

---

## 🔧 Configuration

### Add Allowed Documentation Files

To allow Claude to write to additional root-level documentation files, edit:

**File:** `.github/scripts/validate-docs-only.sh`

```bash
# Add your file to this condition:
if [[ $file == "README.md" ]] || \
   [[ $file == "API.md" ]] || \
   [[ $file == "ENDPOINTS.md" ]] || \
   [[ $file == "CHANGELOG.md" ]] || \
   [[ $file == "YOUR_NEW_FILE.md" ]]; then  # ← Add here
    echo "  ✅ $file (allowed documentation file)"
    continue
fi
```

### Remove Allowed Files

Simply remove the line from the condition above.

### Allow Subdirectories in docs/

The current setup allows **ANY** file in `docs/`:
```bash
if [[ $file == docs/* ]]; then
    echo "  ✅ $file (docs folder)"
    continue
fi
```

This is intentional - docs/ is the designated documentation folder.

---

## 🧪 Testing the Validation

### Test 1: Valid Documentation Commit

```bash
# Manually test the validation script
echo "test" > docs/test.md
git add docs/test.md
.github/scripts/validate-docs-only.sh
# Should output: ✅ All staged files are valid
```

### Test 2: Invalid Code Commit

```bash
# Try to add a source file
git add internal/handler/pokemon.go
.github/scripts/validate-docs-only.sh
# Should output: ❌ ERROR and fail
```

---

## 🚨 Security Considerations

### Why This Matters

1. **Prevents Accidental Code Changes**
   - Claude won't accidentally modify source code while updating docs
   - Reduces risk of introducing bugs

2. **Separation of Concerns**
   - Code changes reviewed by developers
   - Documentation changes automated by Claude
   - Clear boundary between the two

3. **Audit Trail**
   - Easy to see what Claude modified (docs only)
   - Source code changes always from developers
   - Better traceability

4. **Compliance**
   - Some organizations require human review of code
   - This ensures Claude only touches documentation
   - Meets security policy requirements

### Limitations

1. **GitHub Actions Permissions**
   - The workflow still has `contents: write` permission
   - This is repository-wide (GitHub doesn't support folder-level)
   - Protection is enforced at the git hook level, not GitHub API level

2. **Bypass Possibility**
   - A malicious actor could modify the workflow to remove the hook
   - Requires write access to `.github/workflows/` folder
   - Mitigated by: branch protection rules and workflow approval

3. **Manual Validation**
   - Developers should still review all commits from Claude
   - Even documentation changes should be verified
   - Trust but verify

---

## 🔐 Best Practices

### 1. Enable Branch Protection

In GitHub repository settings:
```
Settings → Branches → Branch protection rules → Add rule

Rule: main (or master)
✅ Require a pull request before merging
✅ Require status checks to pass
✅ Include administrators
```

### 2. Review Claude's Commits

Always review documentation commits:
- Check the diff in the PR
- Verify Postman collection is accurate
- Test the updated documentation
- Approve before merging

### 3. Monitor Workflow Changes

Set up notifications for changes to `.github/workflows/`:
- These are critical security files
- Should be reviewed by senior developers
- Consider requiring CODEOWNERS approval

### 4. Regular Audits

Periodically review:
- What files Claude has committed
- Whether validation is working correctly
- If any suspicious activity occurred

---

## 📊 Audit Log

### Where to Check Claude's Activity

1. **PR Commits Tab**
   - Shows all commits in a PR
   - Author will be "Claude Code Bot"
   - Message will start with "docs:"

2. **GitHub Actions Logs**
   - Go to Actions tab
   - Find the "Claude Code" workflow
   - Check the logs for validation output

3. **Git History**
   ```bash
   # See all commits from Claude
   git log --author="Claude Code Bot" --oneline

   # See files Claude modified
   git log --author="Claude Code Bot" --name-only
   ```

---

## ❓ FAQ

### Q: Can I give Claude write access to other folders?

**A:** Yes, edit `.github/scripts/validate-docs-only.sh` to allow additional patterns:

```bash
# Example: Allow writing to api/ folder
if [[ $file == docs/* ]] || [[ $file == api/* ]]; then
    echo "  ✅ $file (allowed folder)"
    continue
fi
```

### Q: What if I need Claude to update go.mod?

**A:** Don't. Dependencies should be managed by developers, not automation. If Claude suggests a dependency update, do it manually.

### Q: Can Claude still read source code?

**A:** Yes! Claude has full READ access. It can analyze code, review PRs, run tests, and provide feedback. It just can't WRITE source code.

### Q: What happens if the validation script fails?

**A:** The commit is blocked. Claude will see an error message and should report it. The PR will not be updated until the issue is fixed.

### Q: Can I disable this validation temporarily?

**A:** Not recommended, but you can:
```bash
# Remove the hook
rm .git/hooks/pre-commit
```

Remember to re-enable it afterwards!

---

## 🆘 Troubleshooting

### Validation script not running

**Check:**
```bash
# Is the hook installed?
ls -la .git/hooks/pre-commit

# Is it executable?
chmod +x .git/hooks/pre-commit
chmod +x .github/scripts/validate-docs-only.sh
```

### Claude's commit still succeeded with code files

**Check:**
1. Did the workflow install the hook?
2. Check GitHub Actions logs
3. Verify the script is correctly blocking

### False positives (valid docs blocked)

**Check:**
- Is the file actually in `docs/` folder?
- Is it one of the allowed root files?
- Check for typos in file path
- Add the file to allowed list if needed

---

**Last Updated:** 2026-02-05
**Maintained By:** Development Team
**Security Level:** Medium (Automated documentation only)
