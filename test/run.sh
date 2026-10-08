#!/bin/bash
# Golden tests: each row runs a command and diffs its stdout + stderr + exit code against test/golden/<name>.
#   test/run.sh [offline]  pinned clock, no network: local stand-ins replace OSM
#   test/run.sh live       real OpenStreetMap and Apple lookups
# UPDATE=1 rewrites goldens. Rows override the environment per call (e.g. alfred_workflow_data=<fixture dir>).
set -u
cd "$(dirname "$0")/.."

tmp=$(mktemp -d)
trap 'kill $(jobs -p) 2>/dev/null; rm -rf "$tmp"' EXIT
export TIMEIN_NOW=2025-05-11T18:38:07Z # Mon, May 12, 1:38 AM in Bangkok
export TIMEIN_LOOKUPS= alfred_workflow_data=$tmp/empty
fx=test/fixtures
refused=http://127.0.0.1:1/search?q= # nothing listens on port 1
rows=0 failed=0

check() { # <name> <command...>
  local name=$1 got
  shift
  rows=$((rows + 1))
  got=$("$@" 2>&1; echo "exit $?")
  if [[ ${UPDATE:-} ]]; then echo "$got" > "test/golden/$name"; return; fi
  diff -u "test/golden/$name" - <<< "$got" || { echo "FAIL: $name"; failed=$((failed + 1)); }
}

# How the script gets run: from a shell, from another directory, and as Alfred runs the installed package.
timein() { osascript -l JavaScript timein.js "$@"; }
from() { (cd "$1" && shift && "$@"); }
alfred() {
  (cd "$tmp/pkg" && bash -c "$(osascript -l JavaScript -e 'ObjC.deepUnwrap($.NSDictionary.dictionaryWithContentsOfFile("info.plist"))
    .objects.find(o => o.type === "alfred.workflow.input.scriptfilter").config.script')" -- "$@")
}

# <status> <body> <command...>: runs the command against a one-shot local Nominatim stand-in.
# Status "silent" accepts the connection and never answers.
nominatim() {
  local port=$((20000 + RANDOM % 20000)) status=$1 body=$2
  shift 2
  if [[ $status == silent ]]; then
    sleep 10
  else
    printf 'HTTP/1.1 %s\r\nContent-Type: application/json\r\nContent-Length: %d\r\nConnection: close\r\n\r\n%s' \
      "$status" "$(($(printf %s "$body" | wc -c)))" "$body"
  fi | nc -l 127.0.0.1 "$port" >/dev/null 2>&1 &
  for _ in {1..100}; do lsof -nP -iTCP:"$port" -sTCP:LISTEN >/dev/null 2>&1 && break; sleep 0.05; done
  TIMEIN_NOMINATIM="http://127.0.0.1:$port/search?q=" "$@"
}

offline() {
  check zone-alfred         timein --format=alfred Asia/Bangkok
  check zone-plain          timein America/New_York
  check zone-case           timein asia/tokyo
  check zone-half-hour      timein --format=alfred America/St_Johns
  check zone-45-minutes     timein --format=alfred Asia/Kathmandu
  check zone-utc            timein --format=alfred GMT
  check zone-ahead-a-day    timein Pacific/Kiritimati
  check zone-behind-a-day   timein Pacific/Pago_Pago
  TIMEIN_NOW=2025-01-15T12:00:00Z check zone-winter timein --format=alfred America/New_York
  TIMEIN_NOW=2025-05-11T17:00:00Z check midnight    timein Asia/Bangkok
  TIMEIN_NOW=2025-05-11T05:00:00Z check noon        timein Asia/Bangkok
  check seed-hit            timein --format=alfred "  PARIS  "
  check seed-multiword      timein buenos aires
  alfred_workflow_data=$fx/override check cache-over-seed    timein london
  alfred_workflow_data=$fx/invalid  check cache-unknown-zone timein bangkok
  alfred_workflow_data=$fx/corrupt  check cache-corrupt      timein tokyo
  check no-match-alfred     timein --format=alfred NotARealCity123
  check no-match-plain      timein NotARealCity123
  check prototype-key       timein constructor # maps are prototype-free: no Object.prototype members
  check proto-key           timein __proto__
  check empty-alfred        timein --format=alfred ""
  check empty-plain         timein ""
  check unknown-format      timein --format=constructor paris

  # Entry points: the shebang from elsewhere (seed found next to the script), the CLI's own cache dir,
  # and the packaged workflow — quarantined as if downloaded — run through its Script Filter.
  check shebang-elsewhere   from "$tmp" "$PWD/timein.js" bangkok
  HOME=$PWD/$fx/home check default-cache-dir env -u alfred_workflow_data "$PWD/timein.js" machu picchu
  make -s alfredworkflow OUT="$tmp/TimeIn.alfredworkflow" >/dev/null
  unzip -q "$tmp/TimeIn.alfredworkflow" -d "$tmp/pkg"
  xattr -w com.apple.quarantine "0083;$(printf %x "$(date +%s)");Safari;" "$tmp"/pkg/*
  check package-contents    unzip -Z1 "$tmp/TimeIn.alfredworkflow"
  check alfred-seed-hit     alfred tokyo
  check alfred-empty        alfred ""

  # OSM failures end in a clean "Could not geocode" — no crash, no hang.
  TIMEIN_LOOKUPS=osm TIMEIN_NOMINATIM=$refused check osm-refused timein "Eiffel Tower"
  TIMEIN_LOOKUPS=osm nominatim '400 Bad Request' "$(<$fx/nominatim/error-400.json)" check osm-http-error timein "Eiffel Tower"
  TIMEIN_LOOKUPS=osm nominatim '200 OK' '<html>Service unavailable</html>' check osm-not-json timein --format=alfred "Eiffel Tower"
  TIMEIN_LOOKUPS=osm nominatim silent '' check osm-timeout timein "Eiffel Tower" # waits out the 5 s network timeout
  TIMEIN_LOOKUPS=osm,mapkit nominatim '200 OK' "$(<$fx/nominatim/no-results.json)" check osm-no-results timein "Eiffel Tower" # MapKit not asked
}

live() {
  export alfred_workflow_data=$tmp/live
  paced() { sleep 1; check "$@"; } # Nominatim allows 1 req/s
  TIMEIN_LOOKUPS=mapkit alfred_workflow_data=$tmp/mapkit check live-mapkit timein Sydney Opera House
  TIMEIN_LOOKUPS=osm,mapkit TIMEIN_NOMINATIM=$refused alfred_workflow_data=$tmp/osm-down \
    check live-osm-down timein Sydney Opera House # OSM unreachable → MapKit answers
  TIMEIN_LOOKUPS=osm alfred_workflow_data=$tmp/stub nominatim '200 OK' "$(<$fx/nominatim/eiffel-tower.json)" \
    check live-apple-zone timein Eiffel Tower # captured OSM reply → Apple's reverse geocoder
  TIMEIN_LOOKUPS=osm        paced live-osm-landmark timein Eiffel Tower
  TIMEIN_LOOKUPS=osm        paced live-osm-airport  timein JFK
  TIMEIN_LOOKUPS=osm,mapkit paced live-not-found    timein NotARealCity123 # OSM's "no such place" must not reach MapKit
  TIMEIN_LOOKUPS=osm,mapkit paced live-no-zone      timein South Pole      # located, but Apple has no zone there
  check live-cached timein eiffel tower # offline: must come from the cache written above
}

case ${1:-offline} in
  offline | live) "${1:-offline}" ;;
  *) echo "usage: $0 [offline|live]" >&2; exit 2 ;;
esac
echo "$((rows - failed))/$rows passed"
exit $((failed > 0))
