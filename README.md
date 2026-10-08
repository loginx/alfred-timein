# alfred-timein

**Fast timezone lookup and time conversion for global collaboration**

An Alfred workflow for instantly finding the current time in any city worldwide. Built for remote teams, frequent travelers, and anyone coordinating across time zones.

It ships no binaries: the whole workflow is one script run by macOS's own JavaScript runtime (JXA), so there is nothing for Gatekeeper to quarantine.

## Quick Start

<p align="center">
  <img src="workflow/screenshot.png" alt="Alfred Timein Workflow Screenshot" width="600" />
</p>

Type in Alfred:

```bash
timein bangkok
timein new york
timein eiffel tower
```

And get:

```text
Asia/Bangkok - Mon, May 12, 1:38 AM
Current time in Bangkok (GMT+7)
```

- ⏎ Copy the full time string
- ⌘⏎ Copy the timezone name (e.g. `Asia/Bangkok`)
- ⌥⏎ Copy ISO 8601 time
- ⌘L Large Type

Or use the script from a terminal:

```bash
./timein.js bangkok
Asia/Bangkok - Mon, May 12, 1:38 AM

./timein.js America/New_York     # IANA zone names skip geocoding entirely
America/New_York - Sun, May 11, 2:38 PM

./timein.js --format=alfred "Eiffel Tower"
{"items":[{"title":"Europe/Paris - Sun, May 11, 8:38 PM","subtitle":"Current time in Paris (GMT+2)", ...}]}
```

## What it understands

- **Cities**: `London` → `Europe/London`
- **Landmarks**: `Eiffel Tower` → `Europe/Paris`
- **Airports**: `JFK` → `America/New_York`
- **Postal codes** (add the country): `90210 USA` → `America/Los_Angeles`, `SW1A 1AA` → `Europe/London`
- **IANA zones**, any case: `asia/tokyo` → `Asia/Tokyo`

## Installation

1. Download the latest release from [the Releases page](https://github.com/loginx/alfred-timein/releases/latest).
2. Double-click `TimeIn.alfredworkflow` to install it in Alfred.
3. In Alfred, type `timein berlin`.

To build it yourself: clone the repo and run `make alfredworkflow`. There is no compile step.

## How it works

1. **IANA zone typed?** Use it as-is.
2. **Cache, then seed**: your past lookups, then the capitals shipped in `capitals.json`.
3. **OpenStreetMap** (Nominatim) turns the query into coordinates. No API key, no location bias.
4. **Apple's geocoder** (`CLGeocoder`) turns those coordinates into a timezone.
5. **MapKit search** is the fallback when OpenStreetMap is unreachable or a future macOS drops `CLGeocoder`.

Results from steps 3–5 are cached. Times and abbreviations come from macOS's own timezone database, so they stay current with OS updates.

| Path | Typical time |
|------|--------------|
| IANA zone, cache or seed hit | ~70 ms |
| First lookup of a place | ~0.6 s |

Design rationale: [docs/architecture](docs/architecture/README.md).

## Caching

- Inside Alfred, lookups are cached in the workflow's data folder (`cache.json`). From a terminal, they go to `~/Library/Caches/com.loginx.timein/cache.json`.
- Your cache wins over the shipped seed, so a wrong answer can be fixed by editing `cache.json`.
- Delete `cache.json` to start over. A corrupt file is treated as empty.
- An entry naming a zone your macOS doesn't know (e.g. `Europe/Kyiv` before tzdata 2022b) is looked up again.

## Compatibility

The script uses ES2020 syntax (JavaScriptCore from Safari 13.1+) and macOS frameworks that date back to 10.11. That puts the floor at macOS 10.15.4, or an older macOS with Safari 13.1 or later installed. Only macOS 27 is tested.

`CLGeocoder` is deprecated as of macOS 26. Its replacement, `MKReverseGeocodingRequest`, never calls back from JXA, so the script keeps using `CLGeocoder` while it exists and falls back to MapKit search if it disappears.

## Testing

Unit tests (`test/unit.js`) pin down the lookup chain; golden tests (`test/run.sh`) run the script as a shell or Alfred would and diff stdout, stderr and exit code against `test/golden/<name>`. Commands are in [CONTRIBUTING.md](CONTRIBUTING.md#testing).

## Known Limitations

- The first lookup of a place that isn't cached or seeded needs an internet connection.
- Places Apple has no timezone for, such as summits and the poles, return an error rather than a guess.
- Abbreviations are the ones macOS uses: `EDT`, but `GMT+7` rather than `ICT`.
- Dates and times are formatted in English.

## License

MIT. Feel free to fork and improve.

Made for Alfred users and CLI fans who prefer speed, simplicity, and control.
