## Case-Insensitive Header Comparison   
Header names comparison is normally case-sensitive.  
To make this comparison case-insensitive, add the `--case-insensitive-headers` flag:
```
oasdiff diff data/header-case/base.yaml data/header-case/revision.yaml --case-insensitive-headers
```

You can ignore multiple elements with a comma-separated list of excluded elements as in [this example](#ignore-changes-to-description-and-examples).  

