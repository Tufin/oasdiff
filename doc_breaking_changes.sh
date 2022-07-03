#!/bin/sh

echo "## Examples of breaking changes"
grep -h 'is .*breaking$' diff/*.go | sed 's/^.*\/\/ /- /g' | sort

echo ""
echo "## Examples of non-breaking changes"
grep -h 'isn.*breaking$' diff/*.go | sed 's/^.*\/\/ /- /g' | sort