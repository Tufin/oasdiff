#!/bin/sh

if [[ $1 == "-v" ]]; then
    echo "Test IDs with line numbers:"
    grep -nRE --exclude-dir=localizations,localizations_src "// (CL|BC): .*: .*-.*" checker/* | cut -d":" -f1,2,5
else
    cmd='grep -nRE --exclude-dir=localizations,localizations_src "// (CL|BC): .*: .*-.*" checker/* | cut -d":" -f5 | tr , '\n' | sort -u'
    cmd_with_wc=$cmd"|wc -l|xargs echo"
    count=$(eval "$cmd_with_wc")
    echo "Unique test IDs ("$count"):"
    eval "$cmd"
fi
