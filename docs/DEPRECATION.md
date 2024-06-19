## Deprecating APIs
Sometimes APIs need to be removed, for example, when we replace an old API by a new version.  
As API owners, we want a process that will allow us to phase out the old API version and transition to the new one smoothly as possible and with minimal disruptions to business.

OpenAPI specification supports a ```deprecated``` flag which can be used to mark operations and other object types as deprecated.  
Normally, deprecation **is not** considered a breaking change since it doesn't break the client but only serves as an indication of an intent to remove something in the future, in contrast, the eventual removal of a resource **is** considered a breaking change.

### Deprecation without a sunset date
Oasdiff allows you to gracefully remove a resource without getting a breaking change error, as follows:
1. The resource is marked as ```deprecated```:
   ```
   /api/test:
    get:
     deprecated: true
   ```
2. Subsequently, the resource can be removed without triggering a breaking change error.

### Deprecation with a sunset date
A more mature deprecation process includes a sunset date which tells the API clients how long they can still use the deprecated API before sunset:
1. The resource is marked as ```deprecated``` and a [special extension](https://swagger.io/specification/#specification-extensions) ```x-sunset``` is added to announce the date at which the resource will be removed:
   ```
   /api/test:
    get:
     deprecated: true
     x-sunset: "2022-12-31"
   ```
   Note sunset date should conform with [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339).    
2. At the sunset date or anytime later, the resource can be removed without triggering a breaking change error. An earlier removal will be considered a breaking change.

Notes:
1. In this mode, oasdiff considers the `x-sunset` extension as **optional** allowing developers to deprecate APIs with or without the `x-sunset` extension.  
2. After an `x-sunset` extension is specified, changing it to an earlier date will trigger a breaking change.  

### Enforcing Sunset Grace Period
Oasdiff allows you to enforce a sunset grace period: the minimal number of days required between deprecating a resource and removing it.  
In this mode oasdiff considers the `x-sunset` extension **mandatory**, so that deprecating an API must be acompanied with an `x-sunset` extension, otherwise it is considered breaking.  

For example, the following command requires deprecation of APIs to be accompanied by an ```x-sunset``` extension with a date which is at least 180 days away, otherwise the deprecation itself will be considered a breaking change:
```
oasdiff breaking data/deprecation/base.yaml data/deprecation/deprecated-past.yaml --deprecation-days-stable=180
```

If you are using [API stability levels](STABILITY.md) you can define different sunset grace periods for:
- Stable APIs: with the `--deprecation-days-stable`
- Beta APIs: with the `--deprecation-days-beta`

Notes:
1. Deprecation days can be set to non-negative integers
2. Setting deprecation days to a zero value disables enforcement and reverts to the [Deprecation with a sunset date](#deprecation-with-a-sunset-date) behavior
2. After an `x-sunset` extension is specified, it can only be changed to a future date which respects the sunset grace period relative to date of the change.