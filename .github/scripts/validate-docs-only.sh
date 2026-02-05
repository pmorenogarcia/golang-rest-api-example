#!/bin/bash
# Validates that only docs/ folder files are being committed
# This ensures Claude can only write to documentation

set -e

echo "🔍 Validating that only documentation files are being committed..."

# Get list of staged files
STAGED_FILES=$(git diff --cached --name-only)

if [ -z "$STAGED_FILES" ]; then
    echo "✅ No files staged for commit"
    exit 0
fi

echo "Staged files:"
echo "$STAGED_FILES"

# Check each file
INVALID_FILES=()
while IFS= read -r file; do
    # Allow files in docs/ folder
    if [[ $file == docs/* ]]; then
        echo "  ✅ $file (docs folder)"
        continue
    fi

    # Allow specific root-level documentation files
    if [[ $file == "README.md" ]] || \
       [[ $file == "API.md" ]] || \
       [[ $file == "ENDPOINTS.md" ]] || \
       [[ $file == "CHANGELOG.md" ]]; then
        echo "  ✅ $file (allowed documentation file)"
        continue
    fi

    # Any other file is not allowed
    echo "  ❌ $file (NOT in docs folder or allowed list)"
    INVALID_FILES+=("$file")
done <<< "$STAGED_FILES"

# If there are invalid files, fail
if [ ${#INVALID_FILES[@]} -gt 0 ]; then
    echo ""
    echo "❌ ERROR: Claude can only commit to the docs/ folder!"
    echo ""
    echo "Invalid files detected:"
    printf '  - %s\n' "${INVALID_FILES[@]}"
    echo ""
    echo "Allowed locations:"
    echo "  ✅ docs/**/*"
    echo "  ✅ README.md, API.md, ENDPOINTS.md, CHANGELOG.md"
    echo ""
    echo "Please remove non-documentation files from this commit."
    exit 1
fi

echo ""
echo "✅ All staged files are valid documentation files"
exit 0
