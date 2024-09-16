/*
Package generator generates the breaking-changes and changelog messages for the checker package.
The output, messages.yaml can be used by the checker package instead of the hardcoded messages under localizations_src.
Advatages over manuallly writing the messages:
- The generated ids and messages are consistent according to the logic in the generator.
- The generator can be easily extended to support more messages.
Additional work needed before using the generator:
- Check that all messages are covered by the generator.
- Decide what to do with Russian messages.
*/
package generator
