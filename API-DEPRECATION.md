## Deprecating APIs
Sometimes APIs need to be removed, for example, when we replace an old API by a new version.
As API owners, we want a process that will allow us to phase out the old API version and transition to the new one smoothly as possible and with minimal disruptions to business.

OpenAPI specification supports a ```deprecated``` flag which can be used to mark operations and other object types as deprecated.  
Normally, deprecation **is not** considered a breaking change since it doesn't break the client but only serves as an indication of an intent to remove something in the future, in contrast, the eventual removal of a resource **is** considered a breaking change.

oasdiff allows you to gracefully remove a resource without getting the ```breaking-change``` warning, as follows:
1. First, the resource is marked as ```deprecated``` and a [special extension](https://swagger.io/specification/#specification-extensions) ```x-sunset``` is added to announce the date at which the resource will be removed
   ```
   /api/test:
    get:
     deprecated: true
     x-sunset: "2022-08-10"
   ```
2. At the sunset date or anytime later, the resource can be removed without triggering a ```breaking-change``` warning. An earlier removal will be considered a breaking change.

In addition, oasdiff also allows you to control the minimal number of days required between deprecating a resource and removing it with the `--deprecation-days-beta` and `--deprecation-days-stable` flags, specifying the deprecation days for each [API stability level](https://github.com/Tufin/oasdiff/blob/main/BREAKING-CHANGES.md#api-stability-levels).  
For example, the following command requires any deprecation to be accompanied by an ```x-sunset``` extension with a date which is at least 30 days away, otherwise the deprecation itself will be considered a breaking change:
```
oasdiff breaking data/deprecation/base.yaml data/deprecation/deprecated-past.yaml --deprecation-days-stable=30
```

By default, `--deprecation-days-beta` and `--deprecation-days-stable` are set to 31 and 180, respectively. Setting deprecation-days to 0 allows for non-breaking deprecation, regardless of the sunset date.  
