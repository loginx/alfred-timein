# Architecture

alfred-timein answers one question — *what time is it there?* — for a city, landmark, airport, postal code or IANA zone, in Alfred or a terminal.

## Shape

One JXA script, `timein.js`, run by `/usr/bin/osascript`. No build step, no dependencies, no binary ([ADR-002](ADR-002-native-jxa-runtime.md)).

```
query ─▶ IANA zone? ─▶ cache.json ─▶ capitals.json ─▶ LOOKUPS ─▶ clock ─▶ FORMATS
          (as typed)    (user)        (shipped seed)    osm        NSDateFormatter   plain | alfred
                                                        mapkit
```

The script reads top-down in the same order:

| Landmark | Role |
|----------|------|
| `LOOKUPS` | Ordered network lookups. Each returns a zone name, `null` (couldn't answer → try the next) or `NOT_FOUND` (answered "nowhere" → stop). |
| `FORMATS` | Output table: `ok` and `error` per format. Unknown format fails loudly. |
| `resolve` | Zone as typed → user cache → seed → lookups; caches network answers. |
| `osm` | Nominatim geocode, then `CLGeocoder` reverse geocode of the coordinates. |
| `mapkit` | `MKLocalSearch` text query; fallback only. |
| `clock` | Title, abbreviation and ISO 8601 time from `NSDateFormatter` (`en_US_POSIX`, stable across releases) and `NSISO8601DateFormatter`. |
| `settle` | Spins the run loop until a completion handler fires — JXA has no `await`. |

## Decisions

- **Two-step lookup (OSM, then Apple).** OpenStreetMap geocodes landmarks without biasing results toward the user's location; Apple's reverse geocoder turns coordinates into a zone without shipping polygon data.
- **MapKit only when OSM can't answer.** MapKit's text search finds *something* for almost any input and leans toward the user's region ("big ben" → a street in South Carolina). Asking it after OSM said "no such place" would turn typos into confident wrong answers, so `NOT_FOUND` ends the chain.
- **Prototype-free maps.** `FORMATS`, `LOOKUPS` and parsed JSON have no prototype, so input like `constructor` can't match `Object.prototype` members.
- **Validate zones against this macOS.** Cached and seeded names are checked with `NSTimeZone`; unknown names are misses, which keeps older macOS releases working when tzdata renames a zone.
- **Seed is data, not a build step.** `capitals.json` is a plain `{city: zone}` map, edited by hand.
- **Lazy framework imports.** CoreLocation and MapKit load only on the network path, keeping cache hits near `osascript`'s own startup time.

## JXA pitfalls this code works around

- `$.ClassName` resolution depends on framework import order (`$.NSURLSession` is undefined under Foundation alone) → classes are resolved with `NSClassFromString`.
- 64-bit ObjC integers (`count`, `statusCode`) arrive as strings, and `"0"` is truthy → avoided, or coerced with `Number()`.
- Nested C structs (e.g. `MKCoordinateRegion`) crash the bridge → no region-bounded MapKit searches.
- `console.log` writes to stderr; script output is the `run` handler's return value.

## Testing

- **Unit** (`test/unit.js`): loads `timein.js` with its network edges stubbed and checks the lookup chain — `null` moves on, `NOT_FOUND` stops, unknown zones are skipped, answers are cached — and how Nominatim replies map onto it.
- **Golden, offline** (`test/run.sh`): runs the script as a shell, the shebang and Alfred (the quarantined package through its Script Filter) would. A pinned clock (`TIMEIN_NOW`), no lookups (`TIMEIN_LOOKUPS=`) and one-shot local Nominatim stand-ins (`TIMEIN_NOMINATIM`) make OSM failures — refused, HTTP error, non-JSON, timeout, no results — deterministic.
- **Golden, live** (`test/run.sh live`): real OSM and Apple lookups, including MapKit taking over when OSM is unreachable.
- **CI**: `make lint` holds the ES2020 ceiling; the offline suite also runs on the oldest supported GitHub-hosted macOS.
