# alfred-timein Architecture

This document describes the Clean Architecture implementation of the alfred-timein project.

## High-Level Architecture

```mermaid
graph TB
    subgraph "User Interfaces"
        CLI[CLI Tools]
        Alfred[Alfred Workflow]
    end

    subgraph "Entry Points"
        GeotzCmd[cmd/geotz]
        TimeinCmd[cmd/timein]
        PreseedCmd[cmd/preseed]
    end

    subgraph "Use Cases Layer"
        GeotzUC[GeotzUseCase]
        TimeinUC[TimeinUseCase]
    end

    subgraph "Domain Layer"
        Location[Location Entity]
        Timezone[Timezone Entity]
    end

    subgraph "Adapters Layer"
        Cache[LRU Cache]
        Geocoder[OpenStreetMap]
        TzFinder[Timezone Finder]
        Presenter[Output Formatters]
    end

    subgraph "External Dependencies"
        OSM[OpenStreetMap API]
        TzDB[Timezone Database]
        FileSystem[Cache Files]
    end

    CLI --> GeotzCmd
    CLI --> TimeinCmd
    Alfred --> GeotzCmd
    Alfred --> TimeinCmd

    GeotzCmd --> GeotzUC
    TimeinCmd --> TimeinUC
    PreseedCmd --> Cache

    GeotzUC --> Location
    GeotzUC --> Timezone
    TimeinUC --> Timezone

    GeotzUC --> Cache
    GeotzUC --> Geocoder
    GeotzUC --> TzFinder
    GeotzUC --> Presenter

    TimeinUC --> Presenter

    Cache --> FileSystem
    Geocoder --> OSM
    TzFinder --> TzDB
```

## Detailed Component Architecture

```mermaid
graph TB
    subgraph "cmd/ - Entry Points"
        GeotzMain[geotz/main.go<br/>CLI argument parsing<br/>Format selection]
        TimeinMain[timein/main.go<br/>Input validation<br/>Output formatting]
        PreseedMain[preseed/main.go<br/>Cache pre-population<br/>Capital cities data]
    end

    subgraph "internal/usecases/ - Application Logic"
        GeotzUC["GeotzUseCase<br/>GetTimezoneFromCity()<br/>Cache-aware lookup"]
        TimeinUC["TimeinUseCase<br/>GetTimezoneInfo()<br/>Time formatting"]
        Interfaces["interfaces.go<br/>Port definitions"]
    end

    subgraph "internal/domain/ - Core Entities"
        LocationEntity[Location<br/>Coordinates<br/>Validation]
        TimezoneEntity[Timezone<br/>IANA timezone<br/>Time operations]
    end

    subgraph "internal/adapters/ - External Interface"
        subgraph "cache/"
            LRUCache[lru.go<br/>TTL-based caching<br/>Persistence to JSON]
        end
        
        subgraph "geocoder/"
            OSMGeocoder[openstreetmap.go<br/>Address to coordinates<br/>HTTP client]
        end
        
        subgraph "timezonefinder/"
            TzfFinder[tzf.go<br/>Coordinates to timezone<br/>Offline database]
        end
        
        subgraph "presenter/"
            AlfredFormatter[alfred.go<br/>JSON for Alfred<br/>Variables & caching]
            PlainFormatter[plain.go<br/>Simple text output]
        end
    end

    subgraph "data/ - Static Data"
        Capitals[capitals.json<br/>46 world capitals<br/>Pre-seed data]
    end

    subgraph "External Systems"
        OSMService[OpenStreetMap API<br/>Geocoding service]
        TzDatabase[Timezone Database<br/>tzf library]
        CacheFiles[geotz_cache.json<br/>Local file system]
    end

    GeotzMain --> GeotzUC
    TimeinMain --> TimeinUC
    PreseedMain --> LRUCache
    PreseedMain --> Capitals

    GeotzUC --> LocationEntity
    GeotzUC --> TimezoneEntity
    TimeinUC --> TimezoneEntity

    GeotzUC --> LRUCache
    GeotzUC --> OSMGeocoder
    GeotzUC --> TzfFinder
    GeotzUC --> AlfredFormatter
    GeotzUC --> PlainFormatter

    TimeinUC --> AlfredFormatter
    TimeinUC --> PlainFormatter

    LRUCache --> CacheFiles
    OSMGeocoder --> OSMService
    TzfFinder --> TzDatabase
```

## Data Flow Diagrams

### Timezone Lookup Flow

```mermaid
sequenceDiagram
    participant User
    participant CLI as CLI/Alfred
    participant UC as GeotzUseCase
    participant Cache
    participant Geocoder
    participant TzFinder
    participant Presenter

    User->>CLI: geotz paris
    CLI->>UC: GetTimezoneFromCity("paris")
    
    UC->>Cache: Get("paris")
    alt Cache Hit
        Cache->>UC: "Europe/Paris"
        UC->>Presenter: FormatTimezone(...)
        Presenter->>CLI: Formatted output
    else Cache Miss
        UC->>Geocoder: GetCoordinates("paris")
        Geocoder->>UC: {lat: 48.8566, lon: 2.3522}
        UC->>TzFinder: GetTimezone(coordinates)
        TzFinder->>UC: "Europe/Paris"
        UC->>Cache: Set("paris", "Europe/Paris")
        UC->>Presenter: FormatTimezone(...)
        Presenter->>CLI: Formatted output
    end
    
    CLI->>User: Display result
```

### Cache Pre-seeding Flow

```mermaid
sequenceDiagram
    participant Preseed as preseed tool
    participant Data as capitals.json
    participant Geocoder
    participant TzFinder
    participant Cache

    Preseed->>Data: Load capital cities
    Data->>Preseed: 46 city entries
    
    loop For each capital
        Preseed->>Geocoder: GetCoordinates(city)
        Geocoder->>Preseed: Coordinates
        Preseed->>TzFinder: GetTimezone(coordinates)
        TzFinder->>Preseed: Timezone
        Preseed->>Cache: Set(city, timezone, 1-year TTL)
    end
    
    Cache->>Preseed: Cache populated
```

## Key Design Principles

### Clean Architecture Layers

1. **Domain Layer** (`internal/domain/`)
   - Core business entities (Location, Timezone)
   - No external dependencies
   - Pure business logic

2. **Use Cases Layer** (`internal/usecases/`)
   - Application-specific business rules
   - Orchestrates domain entities and adapters
   - Defines interfaces for external dependencies

3. **Adapters Layer** (`internal/adapters/`)
   - Implements external interfaces
   - Handles I/O operations
   - Converts between external and internal formats

4. **Frameworks Layer** (`cmd/`)
   - Entry points and CLI parsing
   - Framework-specific code
   - Dependency injection

### Performance Optimizations

```mermaid
graph LR
    subgraph "Cache Strategy"
        PreSeed[Pre-seeded Cache<br/>46 capitals<br/>1-year TTL]
        Runtime[Runtime Cache<br/>User lookups<br/>30-day TTL]
        LRU[LRU Eviction<br/>1000 entries max]
    end

    subgraph "Performance Metrics"
        CacheHit[Cache Hit<br/>5-23ms response]
        CacheMiss[Cache Miss<br/>400ms+ response]
        Speedup[154x speedup ratio]
    end

    PreSeed --> CacheHit
    Runtime --> CacheHit
    LRU --> CacheHit
```

### Testing Architecture

```mermaid
graph TB
    subgraph "Test Types"
        Unit[Unit Tests<br/>Fast, isolated<br/>Domain & adapters]
        Integration[Integration Tests<br/>Real CLI & cache<br/>Performance SLA]
        BDD[BDD Tests<br/>User scenarios<br/>Feature validation]
    end

    subgraph "Test Utilities"
        Preseed[Preseed Regeneration<br/>Clean cache state<br/>Consistent results]
        Mocks[Interface Mocks<br/>Dependency isolation<br/>Fast execution]
    end

    Unit --> Mocks
    Integration --> Preseed
    BDD --> Preseed
```

## File Organization

```text
alfred-timein/
├── cmd/                    # Entry points
│   ├── geotz/             # Timezone lookup CLI
│   ├── timein/            # Time display CLI
│   └── preseed/           # Cache pre-population
├── internal/              # Private application code
│   ├── domain/            # Core business entities
│   ├── usecases/          # Application business rules
│   └── adapters/          # External system interfaces
│       ├── cache/         # LRU cache with persistence
│       ├── geocoder/      # Address to coordinates
│       ├── timezonefinder/# Coordinates to timezone
│       └── presenter/     # Output formatting
├── data/                  # Static data files
├── features/              # BDD test scenarios
├── workflow/              # Alfred workflow assets
└── docs/                  # Documentation
```

This architecture provides:

- **Separation of Concerns**: Each layer has distinct responsibilities
- **Testability**: Interfaces enable easy mocking and testing
- **Performance**: Intelligent caching with pre-seeding
- **Maintainability**: Clear boundaries and minimal coupling
- **Extensibility**: New adapters can be added without changing core logic
