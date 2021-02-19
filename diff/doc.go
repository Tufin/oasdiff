/*
Package diff provides a diff function for OpenAPI Spec 3.
Given two specifications it reports:
  * added, deleted and modified endpoints
  * added, deleted and modified schemas
  * a summary of the changes
The diff can be used as a go struct or a json object.
Work to enhance the diff with additional aspects of OAS is in-progress.
*/
package diff