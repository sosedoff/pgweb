# Contributing

## Pull Requests

1. Everything (within reason) must have BDD-style tests.
2. Test driving (TDD) is very strongly encouraged.
3. Follow all existing patterns and conventions in the codebase.
4. Before issuing a pull-request, please rebase your branch against master.
   If you are okay with the maintainer rebasing your pull request, please say so.
5. After issuing your pull request, check [Travis CI](https://travis-ci.org/sclevine/agouti) to make sure that all tests still pass.

## Development Setup

* Clone the repository.
* Follow the instructions on agouti.org to install Ginkgo, Gomega, PhantomJS, ChromeDriver, and Selenium.
* Run all of the tests using: `ginkgo -r .`
* Start developing!

## Method Naming Conventions

### Agouti package (*Page, *Selection)

These are largely context-dependent, but in general:
* `Name` - Methods that do not have a corresponding getter/setter should not start with "Get", "Is", or "Set".
* `GetName` - Non-boolean methods that get data and have a corresponding `SetName` method should start with "Get".
* `IsName` - Boolean methods that get data and have a corresponding `SetName` method should start with "Is".
* `SetName` - Methods that set data and have a corresponding `GetName` method should start with "Set".
* `ReadName` - Methods that exhaust and return data should start with "Read".
* `EnterName` - Methods that enter data without replacing it should start with "Enter".

### API package (*Session, *Element, *Window)

All API method names should be as close to their endpoint names as possible.
* `GetName` for all GET requests returning a non-boolean
* `IsName` for all GET requests returning a boolean
* `SetName` for POST requests that change the browser state
* `NewNames` for POST requests that return and exhaust some browser state (ex. logs)
* `Name` for POST requests that perform some action or retrieve data
* `GetNameElement` for all POST requests returning an element
