## Delta - a distance function for OpenAPI Spec 3 [Beta]
Delta calculates a numeric value between 0 and 1 representing the distance between base and revision specs.  
For example:
```
oasdiff delta base.yaml revision.yaml
```


### The distance between identical specs is 0
The minimum distance, 0, respresnts the distance between specifications with identical endpoints.  
For example, the distance between any spec to itself is 0:
```
oasdiff delta spec.yaml spec.yaml
```

### The distance between disjoint specs is 1
The maximum distance, 1, respresnts the distance between specifications with no common endpoints.  
For example, the distance between a spec with no endpoints and another spec with one or more endpoints is 1:
```
oasdiff delta empty-spec.yaml non-empty-spec.yaml
```


### Symmetric mode
By default, delta is symmetric and takes into account both elements of base that are deleted in revision and elements of base that are added in revision.  
For example, the symmetric distance between any two specs is the same regardless of the order of the given specs:
```
oasdiff delta base.yaml revision.yaml
oasdiff delta revision.yaml base.yaml
```

### Asymmetric mode
It is also possible to calculate an asymmetric distance which takes into account elements of base that were deleted in revision but ignores elements that are missing in base and were added in revision.  
For examples, for any two specs the sum of the asymmetric distances is their symmetric distance:
```
oasdiff delta base.yaml revision.yaml --asymmetric
oasdiff delta revision.yaml base.yaml --asymmetric
```

### Feature status [Beta]
Delta currently considers:
- Endpoints (path+method)
  - Parameters
    - Schema
      - Type
  - Responses

Other elementes of OpenAPI spec are ignored.  
Please submit feature requests.

