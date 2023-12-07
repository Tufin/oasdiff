# oasdiff-service

### Creating a tenant
Create a tenant and get a tenant ID:
```
curl -d '{"tenant": "my-company", "email": "james@my-company.com"}' https://register.oasdiff.com/tenants
```
You will get a response with your tenant ID:
```
{"id": "2ahh9d6a-2221-41d7-bbc5-a950958345"}
```
### Run diff
```
curl -X POST \
    -F base=@data/openapi-test1.yaml \
    -F revision=@data/openapi-test3.yaml \
    http://api.oasdiff.com/tenants/{tenant-id}/diff
```

### Run breaking-changes
```
curl -X POST \
    -F base=@data/openapi-test1.yaml \
    -F revision=@data/openapi-test3.yaml \
    https://api.oasdiff.com/tenants/{tenant-id}/breaking-changes
```

### Run changelog
```
curl -X POST \
    -F base=@data/openapi-test1.yaml \
    -F revision=@data/openapi-test3.yaml \
    https://api.oasdiff.com/tenants/{tenant-id}/changelog
```
### Errors
oasdiff-service uses conventional HTTP response codes to indicate the success or failure of an API request. In general: Codes in the 2xx range indicate success. Codes in the 4xx range indicate a failure with additional information provided (e.g., invalid OpenAPI spec format, a required parameter was missing, etc.). Codes in the 5xx range indicate an error with oasdiff-service servers (these are rare)
