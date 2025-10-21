#!/bin/bash
# test-tasks.sh - Test all Taskfile targets
set -e

echo "=== Testing Taskfile Targets ==="
echo ""

# Track results
PASSED=0
FAILED=0

test_task() {
    local task_name=$1
    local description=$2
    echo -n "Testing $task_name... "
    if task "$task_name" > /dev/null 2>&1; then
        echo "✓ PASS"
        ((PASSED++))
    else
        echo "✗ FAIL - $description"
        ((FAILED++))
    fi
}

# Info/Config tasks (fast)
echo "--- Info/Config Tasks ---"
test_task "config" "Show configuration"
test_task "list:sdks" "List available SDKs"
test_task "workspace:list" "List workspace modules"
echo ""

# Icon generation (fast)
echo "--- Icon Generation ---"
test_task "icons:hybrid" "Generate icons for hybrid-dashboard"
echo ""

# Summary
echo "=== Test Summary ==="
echo "Passed: $PASSED"
echo "Failed: $FAILED"

if [ $FAILED -eq 0 ]; then
    echo "✓ All tests passed!"
    exit 0
else
    echo "✗ Some tests failed"
    exit 1
fi
