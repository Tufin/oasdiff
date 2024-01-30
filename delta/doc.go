/*
Package delta provides a distance function for OpenAPI Spec 3.
The delta is a numeric value between 0 and 1 representing the distance between base and revision specs.

For any spec, a: delta(a, a) = 0
For any two specs, a and b, with no common elements: delta(a, b) = 1

Symmetric mode:
Delta considers both elements of base that are deleted in revision and elements of base that are added in revision.
For any two specs, a and b: delta(a, b) = delta(b, a)

Asymmetric mode:
Delta only considers elements of base that are deleted in revision but not elements that are missing in base and were added in revision.
For any two specs, a and b:  Delta(a, b) + Delta(b, a) = 1
*/
package delta
