## Configuration Files

Oasdiff can be customized through command-line flags, for example:
```
oasdiff changelog data/openapi-test1.yaml data/openapi-test3.yaml --format yaml
```

To see available flags, run `oasdiff help <cmd>`, for example:
```
oasdiff help changelog
```

The flags can also be provided through a configuration file.  
The config file should be named oasdiff.{json,yaml,yml,toml,hcl} and placed in the directory where the command is run.  
For example, see [oasdiff.yaml](oasdiff.yaml).

Note that some of the flags define paths to additional configuration files:
- --err-ignore string:              configuration file for ignoring errors
- --severity-levels string:         configuration file for custom severity levels
- --warn-ignore string:             configuration file for ignoring warnings

Note that command-line flags take precedence over configuration file settings.