#!/bin/bash

# Test script to validate .gitignore effectiveness
# This script creates test files that should be ignored and verifies they are

set -e

echo "🧪 Testing .gitignore effectiveness..."

# Create temporary test files that should be ignored
TEST_FILES=(
    # Go build artifacts
    "test.exe"
    "test.dll"
    "test.so"
    "test.dylib"
    "test.test"
    "test.out"
    
    # Coverage files
    "coverage.txt"
    "coverage.html"
    "coverage_test.html"
    "coverage.out"
    "test.cover"
    "c.out"
    
    # IDE files
    "test.swp"
    "test.swo"
    "test~"
    
    # OS files
    ".DS_Store"
    "Thumbs.db"
    "Desktop.ini"
    
    # Project specific
    "prettyprint"
    "test.formatted.config"
    "test_sample.config"
    "temp_test.config"
    
    # Logs and temp
    "test.log"
    "test.tmp"
    "test.temp"
    "test.bak"
    "test.backup"
    
    # Security
    ".env"
    "test.key"
    "test.pem"
    
    # Profiling
    "test.prof"
    "test.pprof"
    "cpu.prof"
    "mem.prof"
    "trace.out"
)

# Create test directories that should be ignored
TEST_DIRS=(
    "vendor"
    "tmp"
    ".cache"
    "logs"
    "log"
    ".tmp"
    ".temp"
    "secrets"
    "testdata/output"
    "test-results"
    "coverage"
)

echo "📁 Creating test files and directories..."

# Create test files
for file in "${TEST_FILES[@]}"; do
    echo "test content" > "$file"
done

# Create test directories
for dir in "${TEST_DIRS[@]}"; do
    mkdir -p "$dir"
    echo "test content" > "$dir/test.txt"
done

# Special case: Create .idea directory with files
mkdir -p ".idea"
echo "test" > ".idea/workspace.xml"
echo "test" > ".idea/modules.xml"

# Special case: Create .vscode directory with files
mkdir -p ".vscode"
echo "test" > ".vscode/settings.json"
echo "test" > ".vscode/launch.json"

echo "🔍 Checking which files are ignored..."

FAILED_FILES=()
PASSED_COUNT=0

# Check each test file
for file in "${TEST_FILES[@]}"; do
    if git check-ignore "$file" >/dev/null 2>&1; then
        echo "✅ $file - correctly ignored"
        ((PASSED_COUNT++))
    else
        echo "❌ $file - NOT ignored (should be ignored)"
        FAILED_FILES+=("$file")
    fi
done

# Check each test directory
for dir in "${TEST_DIRS[@]}"; do
    if git check-ignore "$dir" >/dev/null 2>&1 || git check-ignore "$dir/" >/dev/null 2>&1; then
        echo "✅ $dir/ - correctly ignored"
        ((PASSED_COUNT++))
    else
        echo "❌ $dir/ - NOT ignored (should be ignored)"
        FAILED_FILES+=("$dir/")
    fi
done

# Check special directories
if git check-ignore ".idea" >/dev/null 2>&1; then
    echo "✅ .idea/ - correctly ignored"
    ((PASSED_COUNT++))
else
    echo "❌ .idea/ - NOT ignored (should be ignored)"
    FAILED_FILES+=(".idea/")
fi

if git check-ignore ".vscode" >/dev/null 2>&1; then
    echo "✅ .vscode/ - correctly ignored"
    ((PASSED_COUNT++))
else
    echo "❌ .vscode/ - NOT ignored (should be ignored)"
    FAILED_FILES+=(".vscode/")
fi

echo "🧹 Cleaning up test files..."

# Clean up test files
for file in "${TEST_FILES[@]}"; do
    rm -f "$file"
done

# Clean up test directories
for dir in "${TEST_DIRS[@]}"; do
    rm -rf "$dir"
done

# Clean up special directories
rm -rf ".idea"
rm -rf ".vscode"

echo ""
echo "📊 Test Results:"
echo "✅ Passed: $PASSED_COUNT"
echo "❌ Failed: ${#FAILED_FILES[@]}"

if [ ${#FAILED_FILES[@]} -eq 0 ]; then
    echo ""
    echo "🎉 All tests passed! .gitignore is working correctly."
    exit 0
else
    echo ""
    echo "⚠️  Some files are not being ignored:"
    for file in "${FAILED_FILES[@]}"; do
        echo "   - $file"
    done
    echo ""
    echo "💡 Consider adding these patterns to .gitignore if they should be ignored."
    exit 1
fi
