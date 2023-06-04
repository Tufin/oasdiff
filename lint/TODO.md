## lint checks to add
1. WARN - Using default with required parameters or properties, for example, with path parameters. This does not make sense – if a value is required, the client must always send it, and the default value is never used.
2. ERROR - Duplicate endpoints:
    a. equal strings
    b. different path param name (path find)
3. ERROR - Default - The default value represents what would be assumed by the consumer of the input as the value of the schema if one is not provided. Unlike JSON Schema, the value MUST conform to the defined type for the Schema Object defined at the same level. 
For example, if type is string, then default can be "foo" but cannot be 1.
4. ERROR - Schema or content keyword. They are mutually exclusive, see: schema vs content: https://swagger.io/docs/specification/describing-parameters/ 
5. ERROR - duplicate required properties
6. ERROR - duplicate properties
7. ERROR - Bad refs
8. ERROR - yaml/json schema validation
9. ERROR - Enhance Info checks:
   - Info/Contact/URL: The URL pointing to the contact information. This MUST be in the form of a URL.
   - Info/License/URL: The email address of the contact person/organization This MUST be in the form of an email address.
   - Info/Terms Of Service: A URL to the Terms of Service for the API. This MUST be in the form of a URL.
10. ERROR - jsonSchemaDialect: The default value for the $schema keyword within Schema Objects contained within this OAS document. This MUST be in the form of a URI.
11. ERROR - In case a Path Item Object field appears both in the defined object and the referenced object, the behavior is undefined. See the rules for resolving Relative References.