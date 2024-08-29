#!/bin/sh

if [[ $1 == "-v" ]]; then
    grep -nRE --exclude-dir=localizations,localizations_src "// (CL|BC): .*: .*-.*" checker/* | cut -d":" -f1,2,5
else
    grep -nRE --exclude-dir=localizations,localizations_src '// (CL|BC): .*: .*-.*' checker/* | cut -d: -f5 | tr , '\n' | awk '{$1=$1};1' | sort -u
fi
