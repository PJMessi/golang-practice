# GoLang Practice Project
This is just a test application for golang practice.

## Basic Commands
Build application
```
make build
```

Build and run application
```
make run
```

## Directories
`pkg`: Contains general purpose services and utilities.  
`internal`: Contains project specific services and utilities.  
`tests`: Contains integration tests.  
`config`: Contains config package that loads environment variables.  
`cmd`: Separates app's main function into dedicated package that allows us to have multiple entry points if needed. Currently has dedicated package for restapi only.  

## Unit Tests
Unit tests for a package is located in the same directory with with filename of orginal_pkg_filename_unit_test.go.

```
make testunit
```

## Integration Tests
Integration tests are all located in tests directory. Nothing is mocked in integration tests in order to make the test environment close to the prod environment as much as possible. So we need to spin up the actual mysql and NATS services specifically for testing. There is a docker compose file `docker-compose-ci.yml` which is used for integration tests in the GitHub workflow. Same can be used to run the integration tests locally.

Everytime changes is pushed to main branch, github workflow will run both the unit tests and integration tests.

```
make testintegration
```

## Whole Tests
Runs both unit and integration tests
```
make test
```
Runs both unit and integration tests and computes unit tests coverage
```
make testcov
```
Runs both unit and integration tests, computes unit test coverage, and generates a coverage report.
```
make testcovrep
```

## Env Variables
While starting the application, it loads the environment variables from the host machine. If there is a `.env` file in the root directory of the project, the values inside it will take precedence. For integration tests, the values inside `.env.test` will take precedence if it is present in the root directory.

