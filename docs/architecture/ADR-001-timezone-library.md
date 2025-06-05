# ADR-001: Timezone Library Selection

## Status
Accepted

## Context
The application needs to convert geographic coordinates to IANA timezone identifiers. This is a core business requirement for timezone resolution functionality.

## Decision
Use the `ringsaturn/tzf` library for coordinate-to-timezone conversion despite its large binary size impact (~15MB).

## Rationale

### Requirements Analysis
1. **Accuracy**: Must provide correct IANA timezone strings for any coordinate
2. **Reliability**: Must work offline (no network dependencies)  
3. **Performance**: Must handle lookups quickly for real-time Alfred workflow
4. **Maintenance**: Prefer stable, well-maintained libraries

### Options Considered

#### Option 1: tzf Library (Chosen)
- **Pros**: Offline, accurate, fast lookups, well-maintained, complete timezone boundary data
- **Cons**: Large binary size (~15MB added), memory usage for data
- **Verdict**: Reliability and offline capability outweigh size concerns

#### Option 2: Online Timezone API
- **Pros**: Small binary size, always up-to-date data
- **Cons**: Network dependency, latency, rate limits, API key management, failure points
- **Verdict**: Rejected due to reliability concerns for Alfred workflow

#### Option 3: Simplified Timezone Mapping
- **Pros**: Tiny binary size, fast
- **Cons**: Inaccurate for border regions, missing edge cases, maintenance burden
- **Verdict**: Rejected due to accuracy requirements

### Implementation Details
- Library provides coordinate → timezone string conversion
- Embedded timezone boundary data ensures offline operation
- Single function call interface fits Clean Architecture patterns
- Well-tested with extensive geographic coverage

## Consequences

### Positive
- Reliable offline timezone resolution
- Accurate results for any global coordinate
- Fast lookup performance (no network latency)
- No API key management or rate limiting concerns
- Fits well with Alfred's instant-response expectations

### Negative  
- Binary size impact (~15MB) affects download and disk usage
- Memory usage for embedded timezone data
- Larger than typical CLI tools

### Mitigation
- Binary size acceptable for desktop applications
- UPX compression available if size becomes critical (53% reduction)
- Alternative: API-based approach can be implemented later if requirements change

## Notes
This decision prioritizes reliability and user experience over binary optimization. For a timezone lookup tool, accuracy and offline capability are more valuable than a smaller download size.