# Contributing

## Pre-requisites
- Go 1.20 or later
- yq 4.0 or later

## Development Environment Setup

To set up your development environment, run the following command:

```sh
make setup
```

This command will:

1.  Download the necessary Go version.
2.  Install Helm v3.

After running `make setup`, you can run the tests to verify the setup:

```sh
make test
```

## Testing

- Unit Test
  ```sh
  make test
  ```

- Integration Test
  ```sh
  # build the plugin from source and place it in the local Helm plugins directory
  make integ-test
  ```

- Integration Test (Snapshot not matched)
  ```sh
  # This will test fail cases where the snapshot does not match the current output.
  make integ-test-fail
  ```

- Integration Test (Kong)
  ```sh
  # This will run the integration tests against a Kong chart submodule.
  make integ-test-kong
  ```
