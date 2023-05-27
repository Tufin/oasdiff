## lint checks to add
1. Using default with required parameters or properties, for example, with path parameters. This does not make sense – if a value is required, the client must always send it, and the default value is never used.
2. Path params without required: true
3. Bad refs
4. Duplicate endpoints
5. Query param mismatch (name in path is differnt than name under 'parameters')
6. Default - The default value represents what would be assumed by the consumer of the input as the value of the schema if one is not provided. Unlike JSON Schema, the value MUST conform to the defined type for the Schema Object defined at the same level. For example, if type is string, then default can be "foo" but cannot be 1.
7. Schema or content keyword. They are mutually exclusive
8. Query string parameters may only have a name and no value
9. non-existing properties marked as required
