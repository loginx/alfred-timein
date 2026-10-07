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

Black-box golden tests at the CLI seam (`test/run.sh`): a pinned clock (`TIMEIN_NOW`) and no network (`TIMEIN_LOOKUPS=`) make the default suite deterministic and offline. `--live` adds real OSM/Apple lookups, including the degraded paths (not found, no zone, MapKit-only).
