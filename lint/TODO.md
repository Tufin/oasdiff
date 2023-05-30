## lint checks to add
1. WARN - Using default with required parameters or properties, for example, with path parameters. This does not make sense – if a value is required, the client must always send it, and the default value is never used.
2. ERROR - Duplicate endpoints:
    a. equal strings
    b. different path param name (path find)
3. ERROR - Default - The default value represents what would be assumed by the consumer of the input as the value of the schema if one is not provided. Unlike JSON Schema, the value MUST conform to the defined type for the Schema Object defined at the same level. 
For example, if type is string, then default can be "foo" but cannot be 1.
4. ERROR - Schema or content keyword. They are mutually exclusive, see: schema vs content: https://swagger.io/docs/specification/describing-parameters/ 
5. ERROR - non-existing properties as part of required list
6. ERROR - Bad refs
7. ERROR - Schema validation
