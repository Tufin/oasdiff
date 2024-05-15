## Notes for Go Developers

### Embedding oasdiff into your program
The simplest way to get a diff in your go program is:
```go
diff.Get(&diff.Config{}, spec1, spec2)
```

### Advanced Examples
- [diff](https://pkg.go.dev/github.com/tufin/oasdiff/diff#example-Get)
- [breaking changes](https://pkg.go.dev/github.com/tufin/oasdiff/diff#example-GetPathsDiff)


### OpenAPI References
Note that oasdiff expects [OpenAPI References](https://swagger.io/docs/specification/using-ref/) to be resolved.  
References are normally resolved automatically when you load the spec. In other cases you can resolve refs using [Loader.ResolveRefsIn](https://pkg.go.dev/github.com/getkin/kin-openapi/openapi3#Loader.ResolveRefsIn).