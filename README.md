# alfred-timein

**Fast timezone lookup and time conversion for global collaboration**

A reliable Alfred workflow and CLI toolset for instantly finding current times in any city worldwide. Built for remote teams, frequent travelers, and anyone coordinating across time zones.

## Quick Start

<p align="center">
  <img src="workflow/screenshot.png" alt="Alfred Timein Workflow Screenshot" width="600" />
</p>

Type in Alfred:

```bash
# In Alfred:
timein bangkok
timein new york
timein tokyo
```

And get:

```text
Asia/Bangkok - Mon, May 12, 1:38 AM
Current time in Bangkok (ICT)
```

Or use the CLI directly:

```bash
# Get the current time in a timezone
bin/timein Asia/Bangkok
Monday, 12 May 2025, 1:38:07 AM

# Get the current time in Alfred JSON format (for piping)
bin/timein --format=alfred Asia/Bangkok
{"items":[{"title":"Asia/Bangkok - Mon, May 12, 1:44 AM","subtitle":"Current time in Bangkok (ICT)","arg":"Asia/Bangkok - Mon, May 12, 1:44 AM","variables":{"timezone":"Asia/Bangkok"}}],"cache":{"seconds":60}}

# Get the timezone for a city or landmark
bin/geotz "Eiffel Tower"
Europe/Paris

# Get the timezone for a city in Alfred JSON format
bin/geotz --format=alfred "Eiffel Tower"
{"items":[{"title":"Europe/Paris","subtitle":"Eiffel Tower (cached)","arg":"Europe/Paris","variables":{"city":"Eiffel Tower"}}],"cache":{"seconds":604800}}
```

## Core Capabilities

### Timezone Lookup
Transform any location into its IANA timezone identifier:
- **Cities**: `"London" → "Europe/London"`
- **Landmarks**: `"Eiffel Tower" → "Europe/Paris"`
- **Airports**: `"JFK" → "America/New_York"`
- **Postal codes**: `"90210" → "America/Los_Angeles"`

### Current Time Display  
Get human-readable local time for any timezone:
- **Formatted output**: `"Monday, 12 May 2025, 1:38:07 AM"`
- **Multiple formats**: Plain text or Alfred JSON
- **Locale-aware**: Includes day, date, and time with timezone abbreviation

### Performance Features
- **Intelligent caching**: 6ms response for cached locations
- **Offline timezone data**: No API dependencies for timezone resolution
- **OpenStreetMap geocoding**: No API keys required
- **Universal binaries**: Native performance on Intel and Apple Silicon

### Integration Options
- **Alfred workflow**: Type `timein bangkok` for instant results  
- **CLI tools**: `geotz` and `timein` for scripting and automation
- **Pipeline support**: `geotz Bangkok | timein` for complex workflows

## Installation

**Recommended:**

1. Download the latest release from [the Releases page](https://github.com/loginx/alfred-timein/releases/latest).
2. Double-click the `.alfredworkflow` file to install it in Alfred.
3. In Alfred, type:

    ```bash
    timein berlin
    ```

**Advanced/Development:**

If you want to build and run the workflow or CLI tools yourself:

1. Clone this repo:

    ```bash
    git clone https://github.com/loginx/alfred-timein.git
    cd alfred-timein
    ```

2. Build the Go binaries:

    ```bash
    make build
    ```

3. Use the CLI tools directly from `bin/`, or package the workflow:

    ```bash
    make alfredworkflow
    ```

## Caching Details

- The persistent cache is stored in the `./.cache/` directory (ignored by git).
- The cache maps city names (lowercased) to their resolved IANA timezone.
- On first lookup, the workflow queries OpenStreetMap and resolves the timezone; subsequent lookups are instant and do not require network access.
- You can safely delete the `.cache/` directory to clear the cache.

## Architecture

alfred-timein follows Clean Architecture principles with clear separation of core logic and external dependencies:

```
Core Features (What)       Implementation (How)
├── Timezone Resolution  ←  OpenStreetMap Geocoding + tzf Library  
├── Time Display        ←  Go time package + Custom Formatting
├── Intelligent Caching ←  LRU Cache with JSON Persistence
└── Multi-format Output ←  Plain Text + Alfred JSON Presenters
```

### Core Components
- **Domain Layer**: Timezone and Location entities with system rules
- **Use Cases**: Timezone lookup and time display workflows  
- **Adapters**: Geocoding, timezone finding, caching, and output formatting
- **Interfaces**: CLI tools and Alfred workflow integration

### User Scenarios
The `features/` directory contains executable specifications showing exactly what the system does:
- `timezone_lookup.feature` - Core timezone resolution capabilities
- `time_display.feature` - Time formatting and display requirements  
- `alfred_integration.feature` - Alfred workflow behavior
- `cli_workflow.feature` - Command-line interface behavior
- `caching_behavior.feature` - Performance and persistence requirements

## Testing & Quality

This project uses comprehensive testing to ensure reliability:

```bash
# Run all tests
make test-all

# Unit tests only  
make test

# BDD scenarios only
make test-bdd
```

**Testing Strategy:**
- **Unit Tests**: Logic validation for each component
- **BDD Scenarios**: Living documentation of user requirements
- **Integration Tests**: End-to-end workflow validation
- **Performance Tests**: Cache behavior and response time validation

## Known Limitations

- Requires an internet connection for initial geocoding

## License

MIT. Feel free to fork and improve.

Made for Alfred users and CLI fans who prefer speed, simplicity, and control.