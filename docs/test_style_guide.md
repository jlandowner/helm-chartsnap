# Test Writing Style Guide

This project uses the **Ginkgo** and **Gomega** testing frameworks for behavior-driven development (BDD). Below are the key conventions and practices for writing tests:

## Frameworks and Libraries
- **Ginkgo**: Provides BDD-style constructs like `Describe`, `Context`, and `It` for organizing tests.
- **Gomega**: Offers expressive matchers for assertions.

## Test Structure
1. **Describe**: Groups related tests, typically for a function or feature.
2. **Context**: Defines specific scenarios or conditions under which tests are executed.
3. **It**: Represents individual test cases with a clear description of the expected behavior.

### Example:
```go
Describe("FeatureName", func() {
    Context("when condition A is met", func() {
        It("should behave as expected", func() {
            Expect(actual).To(Equal(expected))
        })
    })
})
```

## Parameterized Tests
- Use `DescribeTable` for tests with multiple input-output combinations.
- Define a `struct` for test cases and use `Entry` to specify individual cases.

### Example:
```go
DescribeTable("function behavior",
    func(tc testCase) {
        result := functionUnderTest(tc.input)
        Expect(result).To(Equal(tc.expected))
    },
    Entry("case 1", testCase{input: "A", expected: "B"}),
    Entry("case 2", testCase{input: "X", expected: "Y"}),
)
```

## Snapshot Testing
- Use `MatchSnapShot` to compare outputs against pre-recorded snapshots.

### Example:
```go
Ω(output.String()).To(MatchSnapShot())
```

## Cleanup
- Use `DeferCleanup` to restore global states or clean up resources after tests.

### Example:
```go
DeferCleanup(func() {
    // Cleanup logic
})
```

## Error Handling
- Use `Expect(err).To(HaveOccurred())` for expected errors.
- Use `Expect(err).ShouldNot(HaveOccurred())` for successful operations.

## Test Execution
- Register the test suite using `RegisterFailHandler` and `RunSpecs`.

### Example:
```go
func TestMain(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Main Suite")
}
```

Follow these conventions to ensure consistency and readability in the test suite.
