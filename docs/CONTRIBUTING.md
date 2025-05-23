# Contributing to Helm Chartsnap

We appreciate your contributions and look forward to collaborating with you!

## How to Contribute
- **Report Bugs**: Use the [GitHub Issues](https://github.com/your-repo/issues) to report bugs.
- **Suggest Features**: Open a feature request in the Issues section.
- **Submit Pull Requests**: Follow the guidelines below to submit changes.

---

## Development Setup
1.  Clone the repository:
    ```bash
    git clone https://github.com/your-repo/helm-chartsnap.git
    cd helm-chartsnap
    ```

2.  Install dependencies:
    ```bash
    make install
    ```

3.  Run tests to ensure everything is working:
    ```bash
    (make test && make integ-test && make integ-test-fail && make integ-test-kong) \
        && echo "All tests passed" || echo "Some tests failed"
    ```

---

## Testing
- Run unit tests:
  ```bash
  make test
  ```

- Run integration tests:
  ```bash
  make integ-test
  ```

- Run integration tests with failure:
  ```bash
  make integ-test-fail
  ```

- Run integration tests using Kong chart:
  ```bash
  make integ-test-kong
  ```
