#!/bin/bash
# Golden tests: each row runs timein.js and diffs stdout + stderr + exit code against test/golden/<name>.
# Offline by default (no network lookups); --live adds real OSM/Apple lookups. UPDATE=1 rewrites goldens.
# Rows override the environment per call, e.g. a fixture cache: alfred_workflow_data=test/fixtures/<dir>.
set -u
cd "$(dirname "$0")/.."

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
export TIMEIN_NOW=2025-05-11T18:38:07Z # Mon, May 12, 1:38 AM in Bangkok
export TIMEIN_LOOKUPS= alfred_workflow_data=$tmp/empty
fx=test/fixtures
rows=0 failed=0

check() { # <name> [args...]
  local name=$1 got
  shift
  rows=$((rows + 1))
  got=$(osascript -l JavaScript timein.js "$@" 2>&1; echo "exit $?")
  if [[ ${UPDATE:-} ]]; then echo "$got" > "test/golden/$name"; return; fi
  diff -u "test/golden/$name" - <<< "$got" || { echo "FAIL: $name"; failed=$((failed + 1)); }
}

check zone-alfred         --format=alfred Asia/Bangkok
check zone-plain          America/New_York
check zone-case           asia/tokyo
check zone-half-hour      --format=alfred America/St_Johns
check zone-45-minutes     --format=alfred Asia/Kathmandu
check zone-utc            --format=alfred GMT
check zone-ahead-a-day    Pacific/Kiritimati
check zone-behind-a-day   Pacific/Pago_Pago
TIMEIN_NOW=2025-01-15T12:00:00Z check zone-winter --format=alfred America/New_York
TIMEIN_NOW=2025-05-11T17:00:00Z check midnight    Asia/Bangkok
TIMEIN_NOW=2025-05-11T05:00:00Z check noon        Asia/Bangkok
check seed-hit            --format=alfred "  PARIS  "
check seed-multiword      buenos aires
alfred_workflow_data=$fx/override check cache-over-seed    london
alfred_workflow_data=$fx/invalid  check cache-unknown-zone bangkok
alfred_workflow_data=$fx/corrupt  check cache-corrupt      tokyo
check no-match-alfred     --format=alfred NotARealCity123
check no-match-plain      NotARealCity123
check prototype-key       constructor # maps are prototype-free: no Object.prototype members
check proto-key           __proto__
check empty-alfred        --format=alfred ""
check empty-plain         ""
check unknown-format      --format=constructor paris

if [[ ${1:-} == --live ]]; then
  export alfred_workflow_data=$tmp/live
  live() { sleep 1; check "$@"; } # Nominatim allows 1 req/s
  TIMEIN_LOOKUPS=mapkit alfred_workflow_data=$tmp/mapkit check live-mapkit Sydney Opera House
  TIMEIN_LOOKUPS=osm        live live-osm-landmark Eiffel Tower
  TIMEIN_LOOKUPS=osm        live live-osm-airport  JFK
  TIMEIN_LOOKUPS=osm,mapkit live live-not-found    NotARealCity123 # OSM's "no such place" must not reach MapKit
  TIMEIN_LOOKUPS=osm,mapkit live live-no-zone      South Pole      # located, but Apple has no zone there
  check live-cached eiffel tower # offline: must come from the cache written above
fi

echo "$((rows - failed))/$rows passed"
exit $((failed > 0))
