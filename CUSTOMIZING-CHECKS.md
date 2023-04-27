# How to Add Custom Breaking-Changes Checks

## Unit Tests and Documentation
1. write a unit test for your scenario and add a comment "BC: \<use-case\> is breaking"
2. optionally, add additional unit tests and comment them with "is breaking" or "is not breaking"
3. Update the breaking-changes examples doc:
    ```
    ./scripts/test.sh
    ```
4. make sure that [BREAKING-CHANGES.md](BREAKING-CHANGES.md) contains your comments

## Localized Backwards Compatibility Messages
1. add localized texts under [checker/localizations_src](checker/localizations_src) (you can use Google Translate for Russian)
2. Update [localization source file](checker/localizations/localizations.go):
    ```
    go-localize -input checker/localizations_src -output checker/localizations
    ```   
3. make sure that [checker/localizations/localizations.go](checker/localizations/localizations.go) contains the new messages

## Write the Checker Function
1. create new go file under [checker](checker) and name it by the breaking-change use case
2. create a check func inside the file and name it accordingly
3. add the checker func to defaultChecks
4. verify that all unit tests pass
