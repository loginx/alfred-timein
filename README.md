# Alfred Workflow: `timein <city>`

## Why Go? Why this workflow?

- **Blazing fast**: Native Go binaries, no Node.js or npm bloat.
- **Tiny footprint**: No dependencies, no runtime, no "node_modules" hell.
- **Use it anywhere**: Works in Alfred and as standalone CLI tools.

## Usage

Search for the current local time in any city via the `timein` keyword in Alfred, or use the CLI tools directly.

<p align="center">
  <img src="about/screenshot.png" alt="Alfred Timein Workflow Screenshot" width="600" />
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

## Description

A fast, zero-dependency Alfred workflow and CLI toolset that tells you the current local time in any city using natural language input. Powered by Go for maximum speed and reliability.

## Features

- Natural city name input (`"timein London"`, `"timein São Paulo"`)
- Lookup by airport codes (e.g., `"timein JFK"` for John F. Kennedy International Airport)
- Lookup by postal codes (e.g., `"timein 90210"` for Beverly Hills, California)
- Lookup by landmarks (e.g., `"timein Eiffel Tower"`, `"timein Statue of Liberty"`)
- Instant results via Alfred Script Filter or CLI
- Location resolution via OpenStreetMap (no API keys required)
- Accurate timezone mapping with IANA strings (`America/Toronto`)
- Persistent disk caching for repeat queries (city → timezone) in `./.cache/`
- Native, universal Go binaries for maximum performance
- Fully tested

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

## Testing

Run Go tests with:

```bash
make test
```

Tests cover core logic: geocoding, timezone resolution, formatting, caching, and edge cases.

## Known Limitations

- Requires an internet connection for initial geocoding

## License

MIT. Feel free to fork and improve.

Made for Alfred users and CLI fans who prefer speed, simplicity, and control.