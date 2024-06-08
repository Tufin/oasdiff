## Deprecating APIs
Sometimes APIs need to be removed, for example, when we replace an old API by a new version.  
As API owners, we want a process that will allow us to phase out the old API version and transition to the new one smoothly as possible and with minimal disruptions to business.

OpenAPI specification supports a ```deprecated``` flag which can be used to mark operations and other object types as deprecated.  
Normally, deprecation **is not** considered a breaking change since it doesn't break the client but only serves as an indication of an intent to remove something in the future, in contrast, the eventual removal of a resource **is** considered a breaking change.

### Deprecation without sunset
Oasdiff allows you to gracefully remove a resource without getting a breaking change error, as follows:
1. The resource is marked as ```deprecated```:
   ```
   /api/test:
    get:
     deprecated: true
   ```
2. Subsequently, the resource can be removed without triggering a breaking change error.

### Deprecation with sunset
A more mature deprecation process includes a sunset date which tells the clients how long they can still use the deprecated API:
1. The resource is marked as ```deprecated``` and a [special extension](https://swagger.io/specification/#specification-extensions) ```x-sunset``` is added to announce the date at which the resource will be removed:
   ```
   /api/test:
    get:
     deprecated: true
     x-sunset: "2022-12-31"
   ```
   Note sunset date should conform with [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339).    
2. At the sunset date or anytime later, the resource can be removed without triggering a breaking change error. An earlier removal will be considered a breaking change.

### Grace period
Oasdiff allows you to specify a required grace period: the minimal number of days required between deprecating a resource and removing it.  
You can configure a grace period by setting `--deprecation-days-stable` to any positive number.  

For example, the following command requires deprecation of APIs to be accompanied by an ```x-sunset``` extension with a date which is at least 30 days away, otherwise the deprecation itself will be considered a breaking change:
```
oasdiff breaking data/deprecation/base.yaml data/deprecation/deprecated-past.yaml --deprecation-days-stable=30
```

You can also define a different grace period for [Beta resources](STABILITY.md) with the `--deprecation-days-beta` flag.
