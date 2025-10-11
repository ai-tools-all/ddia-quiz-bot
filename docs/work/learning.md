# Go Development Learning Notes

## Testing Best Practices

### Table-Driven Tests
```go
tests := []struct {
    name       string
    input      Type
    wantResult Type
    wantError  bool
}{
    {name: "descriptive case", input: ..., wantResult: ...},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got, err := Function(tt.input)
        // assertions
    })
}
```

### Test Coverage Strategy
1. **Happy path** - normal operation
2. **Edge cases** - empty inputs, boundaries, nil values
3. **Error cases** - invalid inputs, missing data
4. **State transitions** - switching between different states
5. **Backward compatibility** - multiple formats/versions

### File-Based Tests
```go
tmpFile, err := os.CreateTemp("", "prefix-*.ext")
require.NoError(t, err)
defer os.Remove(tmpFile.Name())
tmpFile.WriteString(content)
tmpFile.Close()
```

### Component Initialization Tests
- Test that components are properly initialized
- Test that state is reset between uses
- Test boundary conditions (out of bounds, empty slices)
- Verify cleanup (nil checks when switching types)

## Code Organization

### When to Use Pointers
- **Pointer receivers**: When mutating struct state
- **Value receivers**: For read-only operations
- **Return pointers**: For large structs or when nil is valid

### Error Handling
- Return errors, don't panic (except in truly exceptional cases)
- Wrap errors with context: `fmt.Errorf("context: %w", err)`
- Check errors immediately after function calls

### Struct Tags
```go
type Model struct {
    Field string `toml:"field" json:"field"`  // Multiple tags
    Skip  string `toml:"-" json:"-"`          // Skip serialization
}
```

## Testing Insights from MCQ Implementation

### Parser Tests
- Test multiple frontmatter formats (TOML, YAML)
- Test option parsing with various formats (-, *, none, spacing)
- Use temp files for markdown parsing tests
- Test extraction functions independently

### State Machine Tests
- Test each transition explicitly
- Verify components are cleaned up when switching
- Test that navigation works correctly
- Test submission changes state as expected

### Integration Tests
- Test full workflow (init → navigate → submit → save)
- Verify persistence with session loading
- Test MCQ-specific fields are saved/loaded correctly

## Code Quality

### Before Committing
1. `go fmt ./...` - format all code
2. Run relevant tests: `go test ./path/... -v`
3. Run all tests if changing shared code
4. Check `git diff` for unintended changes
5. Write descriptive commit messages

### Test Naming
- `Test<FunctionName>` for unit tests
- `Test<Feature><Scenario>` for integration tests
- Use descriptive names: `TestSwitchingBetweenQuestionTypes` not `TestSwitch`

### Assertions
- Use testify/assert for readability
- Add helpful messages: `assert.Equal(t, want, got, "context about what failed")`
- Use `require` when test can't continue on failure

## Common Patterns

### Type Detection Pattern
```go
if m.currentQType == "mcq" && m.mcqComponent != nil {
    // MCQ-specific logic
} else {
    // Default/subjective logic
}
```

### Component Initialization Pattern
```go
func (m *Model) initializeComponent() {
    if m.currentIndex >= len(m.items) {
        return  // Guard clause
    }
    
    item := m.items[m.currentIndex]
    
    switch item.Type {
    case "typeA":
        m.componentA = NewComponentA(item.Data)
    case "typeB":
        m.componentB = NewComponentB(item.Data)
    default:
        m.componentA = NewComponentA(item.Data) // Default
    }
}
```

### Clean State Transitions
When switching between component types, explicitly:
1. Clear previous component (set to nil)
2. Initialize new component fresh
3. Reset related state flags
4. Load any persisted data

## Key Takeaways

1. **Test state transitions explicitly** - Don't just test happy path
2. **Use table-driven tests** - Easy to add cases, easy to read
3. **Format before committing** - `go fmt` is non-negotiable
4. **Test backward compatibility** - Support multiple formats
5. **Boundary conditions matter** - Empty slices, nil pointers, out of bounds
6. **Cleanup is critical** - Use defer, set to nil when switching types
7. **Descriptive test names** - Future you will thank current you
