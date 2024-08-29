#!/bin/bash

if [[ $1 == "lines" ]]; then
    grep -nRE --exclude-dir=localizations,localizations_src "// (CL|BC): .*: .*-.*" checker/* | cut -d":" -f1,2,5
elif [[ $1 == "ids" ]]; then
    grep -nRE --exclude-dir=localizations,localizations_src '// (CL|BC): .*: .*-.*' checker/* | cut -d: -f5 | tr , '\n' | awk '{$1=$1};1' | sort -u
elif [[ $1 == "diff" ]]; then
    diff --side-by-side <(grep -nRE --exclude-dir=localizations,localizations_src '// (CL|BC): .*: .*-.*' checker/* | cut -d: -f5 | tr , '\n' | awk '{$1=$1};1' | sort -u) <(cat checker/localizations_src/en/messages.yaml | sed '/# warnings/q' | cut -d":" -f1 | sort -n)
fi
