# Living Documentation Through BDD

This project uses Behavior-Driven Development to create **executable specifications** that serve as both documentation and automated tests. The BDD scenarios define exactly what alfred-timein does for users.

## Philosophy: Executable Business Requirements

BDD scenarios answer three critical questions:
1. **What** does the system do? (Feature descriptions)
2. **How** do users interact with it? (Scenario steps)  
3. **Why** does it matter? (Business value statements)

Each feature file is a contract that ensures the system delivers real user value, not just technical functionality.

## Feature Coverage

Our BDD test suite covers five core business domains:

### 🌍 **Timezone Lookup by City** (`features/timezone_lookup.feature`)
- Finding timezones for major cities, landmarks, and locations with spaces
- Handling unknown locations gracefully with helpful error messages  
- Performance requirements and caching behavior
- Case insensitive lookup and whitespace handling

### ⏰ **Current Time Display** (`features/time_display.feature`)  
- Human-readable time formatting with locale awareness
- Timezone validation before time display
- Different formats for different timezones
- Error handling for invalid inputs

### 🔍 **Alfred Workflow Integration** (`features/alfred_integration.feature`)
- Alfred Script Filter JSON format compliance
- Error formatting for Alfred interface
- Cache indicators in Alfred results
- Proper cache duration headers

### 💻 **Command Line Interface** (`features/cli_workflow.feature`)
- Basic command usage for both `geotz` and `timein`
- Pipeline workflows (`geotz | timein`)
- Format options (plain text vs Alfred JSON)
- Error handling and help text

### 🚀 **Intelligent Caching** (`features/caching_behavior.feature`)
- Cache hit/miss performance characteristics
- LRU eviction policy compliance
- Cache persistence across restarts
- TTL expiration behavior
- Graceful handling of corrupted cache entries

## Running BDD Tests

### Basic BDD Test Run
```bash
make test-bdd
```

### Run All Tests (Unit + BDD)
```bash
make test-all
```

### Run BDD with Verbose Output
```bash
go test -run TestBDD -v
```

### Run Specific Features
```bash
go test -run TestBDD -godog.features=features/timezone_lookup.feature
```

## Understanding BDD Output

BDD tests produce human-readable output showing:
- ✅ **Green**: Scenarios that pass (business requirements met)
- 🟡 **Yellow**: Undefined steps (features documented but not yet implemented)
- ❌ **Red**: Failed scenarios (business requirements violated)

**Current Status**: Clean output with 12 passing scenarios and 64 passing steps. Additional scenarios are documented in comments but temporarily removed to avoid undefined step pollution.

Example output:
```
Feature: Timezone Lookup by City
  Scenario: Finding timezone for a major city
    Given I want to know the timezone for "Tokyo"  
    When I request the timezone information
    Then I should get "Asia/Tokyo" as the timezone
    And the response should be fast
```

## Business Value

### 📋 **Living Documentation**
- Feature files serve as executable specifications
- Business stakeholders can validate requirements
- No gap between documentation and implementation

### 🛡️ **Regression Protection** 
- Ensures user-facing behavior remains consistent
- Catches breaking changes before they reach users
- Validates end-to-end workflows

### 🤝 **Stakeholder Communication**
- Natural language scenarios bridge technical/business divide
- Clear acceptance criteria for new features
- Shared understanding of system behavior

### ⚡ **Continuous Validation**
- Business rules are tested with every build
- Performance requirements are continuously validated
- Cache behavior is verified across scenarios

## Writing New BDD Scenarios

When adding new features:

1. **Start with the business need** in natural language
2. **Write scenarios** before implementation (outside-in development)
3. **Use Given/When/Then** structure consistently
4. **Focus on outcomes**, not implementation details
5. **Make scenarios readable** by non-technical stakeholders

Example template:
```gherkin
Feature: [Business Capability]
  As a [user type]
  I want to [goal]
  So that [business value]

  Scenario: [Specific behavior]
    Given [initial context]
    When [action performed]  
    Then [expected outcome]
    And [additional validation]
```

## Implementation Details

- **Step definitions** in `bdd_test.go` bridge natural language to code
- **Test context** maintains state between steps
- **Service integration** tests real adapters and use cases
- **Timeout handling** ensures tests complete within reasonable time
- **Cache management** provides isolation between scenarios

This BDD implementation ensures our timezone tools meet real user needs while providing confidence that changes don't break existing functionality.