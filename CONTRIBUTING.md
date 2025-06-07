# Contributing

## Development Environment Setup

To set up your development environment, run the following command:

```bash
make setup
```

This command will:

1.  Install the Go version specified in the `go.mod` file.
2.  Download the necessary Go modules.
3.  Install Helm v3.

After running `make setup`, you can run the tests to verify the setup:

```bash
make test
```
