## Common Parameters

### Common Parameters Definition
Parameters shared by all operations of a path can be defined on the path level instead of the operation level.  
Path-level parameters are inherited by all operations of that path.  
A typical use case are the GET/PUT/PATCH/DELETE operations that manipulate a resource accessed via a path parameter.

### There are two ways to handle common parameters in oasdiff
1. By default, oasdiff compares path parameters and operation parameters separately.
2. The `--flatten-params` merges common parameters from the path level into the operation level before running the diff.

For example, this command outputs two breaking changes:
```
oasdiff changelog data/common-params/params_in_path.yaml data/common-params/params_in_op.yaml
```
Output: 
```
2 changes: 2 error, 0 warning, 0 info
error	[new-request-path-parameter] at data/common-params/params_in_op.yaml
	in API GET /admin/v0/abc/{id}
		added the new path request parameter 'id'

error	[new-required-request-parameter] at data/common-params/params_in_op.yaml
	in API GET /admin/v0/abc/{id}
		added the new required 'header' request parameter 'tenant-id'
```


Adding the `--flatten-params` eliminates the errors:
```
oasdiff changelog data/common-params/params_in_path.yaml data/common-params/params_in_op.yaml --flatten-params
```
