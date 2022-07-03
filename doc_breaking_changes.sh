#!/bin/sh

echo "## Examples of breaking changes"
grep -h '// BC:' diff/*.go | sed 's/^.*\/\/ BC: /- /g' | grep -v "isn't" | sort

echo ""
echo "## Examples of non-breaking changes"
grep -h '// BC:' diff/*.go | sed 's/^.*\/\/ BC: /- /g' | grep "isn't" | sort