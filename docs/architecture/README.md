# Architecture Overview

alfred-timein is designed around **timezone resolution and time display** as core business capabilities. The architecture prioritizes reliability, performance, and maintainability for global time coordination workflows.

## Business Context

**Problem**: Remote teams and global travelers need instant, reliable timezone information and current time lookup without complex setup or API dependencies.

**Solution**: Fast, offline-capable timezone resolution with intelligent caching and multiple interface options (Alfred + CLI).

## Core Business Capabilities

### 1. Timezone Resolution
**Purpose**: Transform human-readable locations into IANA timezone identifiers

**Flow**: Location Input → Geocoding → Coordinate-based Timezone Lookup → IANA String  
**Examples**: `"Tokyo" → "Asia/Tokyo"`, `"Eiffel Tower" → "Europe/Paris"`

### 2. Current Time Display  
**Purpose**: Show human-readable local time for any timezone

**Flow**: Timezone String → Current Time Calculation → Formatted Output  
**Examples**: `"Asia/Tokyo" → "Monday, 12 May 2025, 2:38:07 PM"`

### 3. Intelligent Caching
**Purpose**: Eliminate repeated network requests for timezone resolution

**Flow**: Cache Check → Network Lookup (if miss) → Cache Store → Response  
**Performance**: 6ms cache hits vs 400ms+ network lookups

## Architectural Patterns

### Clean Architecture
```
Domain Models (Business Rules)
    ↑
Use Cases (Application Logic)  
    ↑
Interface Adapters (Data Conversion)
    ↑  
Frameworks & Drivers (External Tools)
```

### Dependency Inversion
- Business logic depends on interfaces, not implementations
- External services (geocoding, timezone data) are abstracted
- Easy to swap implementations without changing core logic

### Command Query Separation
- **Commands**: Cache operations, no return values
- **Queries**: Timezone lookup and time display, read-only

## Technology Decisions

### Go Language
**Rationale**: Fast compilation, single binary output, excellent concurrency, strong standard library for time handling

### tzf Library (15MB)
**Rationale**: Offline timezone resolution prioritized over binary size - reliability beats optimization for desktop workflows

### OpenStreetMap Geocoding  
**Rationale**: No API keys required, good coverage for cities and landmarks, free usage

### LRU Cache with JSON Persistence
**Rationale**: Simple, fast, survives application restarts, easy to inspect and debug

### BDD Testing with Godog
**Rationale**: Executable specifications serve as living documentation of business requirements

## Performance Characteristics

- **Cache Hit**: ~6ms response time
- **Cache Miss**: ~400ms (network dependent)  
- **Binary Size**: ~15MB (acceptable for desktop tools)
- **Memory Usage**: Minimal, cache bounded by LRU eviction

## Integration Points

### Alfred Workflow
- JSON Script Filter format for native Alfred integration
- Instant results with subtitle information
- Action support for copying results

### CLI Tools  
- `geotz`: Location → Timezone resolution
- `timein`: Timezone → Current time display
- Pipeline support: `geotz Bangkok | timein`

## Quality Attributes

1. **Reliability**: Offline timezone data, persistent caching, comprehensive testing
2. **Performance**: Intelligent caching, optimized binary size, fast startup  
3. **Usability**: Natural language input, multiple output formats, clear error messages
4. **Maintainability**: Clean Architecture, BDD scenarios, conventional commits