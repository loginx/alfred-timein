# Alfred Workflow: `timein <city>`

## Usage

Search for the current local time in any city via the `timein` keyword.

<p align="center">
  <img src="docs/timein-screenshot.png" alt="Alfred Timein Workflow Screenshot" width="600" />
</p>

Type:

```bash
# In Alfred:
timein bangkok
timein new york
timein tokyo

# From the terminal:
npm run timein bangkok
```

And get:

```text
Asia/Bangkok (UTC+7)  Fri, May 2, 9:30 AM
```

## Description

A fast, zero-config Alfred workflow that tells you the current local time in any city using natural language input.

## Features

- Natural city name input (`"timein London"`, `"timein São Paulo"`)
- Lookup by airport codes (e.g., `"timein JFK"` for John F. Kennedy International Airport)
- Lookup by postal codes (e.g., `"timein 90210"` for Beverly Hills, California)
- Lookup by landmarks (e.g., `"timein Eiffel Tower"`, `"timein Statue of Liberty"`)
- Instant results via `alfy` Script Filter
- Built-in debouncing via Alfred for efficient input handling
- Location resolution via OpenStreetMap (no API keys required)
- Accurate timezone mapping with IANA strings (`America/Toronto`)
- **Persistent disk caching for repeat queries (city → timezone) in `./.cache/`**
- Fully tested with Vitest
- ESM-only codebase (Node.js 18+)

## Installation

**Recommended:**

1. Download the latest release from [the Releases page](https://github.com/loginx/alfred-timein/releases/latest).
2. Double-click the `.alfredworkflow` file to install it in Alfred.
3. In Alfred, type:

    ```bash
    timein berlin
    ```

**Advanced/Development:**

If you want to build and run the workflow yourself:

1. Clone this repo into your Alfred workflows folder or develop with `alfred-link`:

    ```bash
    git clone https://github.com/your-username/alfred-timein.git
    cd alfred-timein
    npm install
    ```

2. Link it with `alfred-link` (or package manually):

    ```bash
    npx alfred-link
    ```

3. In Alfred, type:

    ```bash
    timein berlin
    ```

## Tech Stack

| Tool         | Purpose                            |
|--------------|------------------------------------|
| `alfy`       | Alfred integration and Script Filter handling |
| `node-geocoder` | City name to coordinates (OpenStreetMap) |
| `tz-lookup` | Coordinates to IANA timezone string |
| `lru-cache` | Persistent LRU disk caching (city → timezone) |
| `vitest`    | Unit testing framework |
| `Intl.DateTimeFormat` | Native date/time formatting |

## Caching Details

- The persistent cache is stored in the `./.cache/` directory (ignored by git).
- The cache maps city names (lowercased) to their resolved IANA timezone.
- On first lookup, the workflow queries OpenStreetMap and resolves the timezone; subsequent lookups are instant and do not require network access.
- You can safely delete the `.cache/` directory to clear the cache.

## Testing

Run tests with:

```bash
npm run test
```

Tests cover core logic: geocoding, timezone resolution, formatting, caching, and edge cases.

## Known Limitations

- Requires an internet connection for initial geocoding
- Results are English-only for now

## Roadmap Ideas

Want to contribute? Here are some next steps:

- Implement internationalization & localization
- Support `timein here` for local resolution
- Support `timein tz-code` (e.g., UTC, GMT, EST)
- Bundle a cache of major locations in the distribution for offline access
- Encourage people to suggest ideas

## License

MIT. Feel free to fork and improve.

Made for Alfred users who prefer speed, simplicity, and control.