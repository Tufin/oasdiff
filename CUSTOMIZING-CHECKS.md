# How to Add Custom Breaking-Changes Checks

## Unit Tests and Documentation
1. Add a unit test for your scenario in one of the test files under [checker](checker) with a comment "BC: \<use-case\> is breaking"
   - add test specs under [data](data) as needed
2. Update the breaking-changes examples doc:
    ```
    ./scripts/test.sh
    ```
3. Make sure that your use case was added to [BREAKING-CHANGES.md](BREAKING-CHANGES.md)

## Localized Backwards Compatibility Messages
1. ASdd localized texts under [checker/localizations_src](checker/localizations_src) (you can use Google Translate for Russian)
2. Update [localization source file](checker/localizations/localizations.go):
    ```
    go-localize -input checker/localizations_src -output checker/localizations
    ```   
3. Make sure that [checker/localizations/localizations.go](checker/localizations/localizations.go) contains the new messages

## Write the Checker Function
1. Create new go file under [checker](checker) and name it by the breaking-change use case
2. Create a check func inside the file and name it accordingly
3. Add the checker func to defaultChecks
4. Verify that all unit tests pass
5. Optionally, add additional unit tests and comment them with "is breaking" or "is not breaking"
