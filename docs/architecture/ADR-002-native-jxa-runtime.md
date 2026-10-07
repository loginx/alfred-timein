# ADR-002: Native JXA Runtime

## Status
Accepted (incubating on `incubator/jxa`). Supersedes [ADR-001](ADR-001-timezone-library.md).

## Context
The workflow shipped two unsigned Go binaries (`geotz` ~22 MB with embedded `tzf` polygons, `timein` ~4 MB). Downloaded workflows carry `com.apple.quarantine`, and Gatekeeper blocks unsigned, un-notarized binaries, so the workflow stripped the attribute itself before every run. The binaries also pulled in transitive dependencies with no runtime purpose (e.g. `mongo-driver/bson` via `tzf` → `paulmach/orb`), generating Dependabot churn.

## Decision
Rewrite the workflow as a single JXA script run by the Apple-signed `/usr/bin/osascript`:

1. Geocode with OpenStreetMap Nominatim (unchanged from the Go version).
2. Resolve coordinates to a zone with `CLGeocoder` reverse geocoding (replaces `tzf`).
3. Fall back to `MKLocalSearch` only when Nominatim is unreachable or `CLGeocoder` is missing.
4. Format times with `NSDateFormatter` (`en_US_POSIX`) and `NSISO8601DateFormatter`; abbreviations from `NSTimeZone` (replaces `go-timezone`).

## Options Considered

| Option | Verdict |
|--------|---------|
| Sign and notarize the Go binaries | Fixes quarantine, but needs a paid Developer ID and notarization in CI; keeps the 26 MB of binaries and the dependency tree. |
| `MKLocalSearch` alone (one call, returns a zone) | Rejected as primary: results are biased toward the user's region ("big ben" → South Carolina) and JXA can't set a search region (nested structs crash the bridge). |
| `CLGeocoder` forward geocoding | Rejected: address-only ("Eiffel Tower" → Tower, Minnesota; "JFK" → no result). |
| `MKGeocodingRequest` / `MKReverseGeocodingRequest` (macOS 26) | The intended replacements for `CLGeocoder`, but their completion handlers never fire under JXA (tested on macOS 27.0.1), despite correct BridgeSupport metadata. |
| Swift script / AppleScriptObjC | Swift needs the Command Line Tools and compiles on every run; AppleScriptObjC can't pass blocks to completion-handler APIs. |

## Consequences

### Positive
- No binary, no quarantine workaround, no build step, no third-party dependencies.
- Faster first lookups (~0.6 s vs ~0.9 s, which was dominated by loading `tzf` polygons); local paths ~70 ms.
- Zone data and abbreviations track macOS updates instead of a vendored snapshot.

### Negative
- Uncached lookups need the network for both steps; `tzf` resolved coordinates offline.
- Places Apple has no zone for (summits, poles) return an error; `tzf` answered from polygons.
- Abbreviations follow macOS (`GMT+7` instead of `ICT`).
- Depends on `CLGeocoder`, deprecated in macOS 26. The MapKit fallback keeps lookups working if it is removed, with lower precision.
- JXA is quirky and effectively frozen; workarounds are documented in [the architecture overview](README.md#jxa-pitfalls-this-code-works-around).
- The Go `geotz`/`timein` CLIs and their pipe are replaced by one script that accepts places and zone names alike.

### Exit conditions
- If `MKReverseGeocodingRequest` starts calling back under JXA, replace `CLGeocoder` with it.
- If JXA is removed from macOS, revisit signed binaries.
