#!/usr/bin/osascript -l JavaScript
// Current local time for a city, landmark or IANA zone — for Alfred and the shell.
// Plain JXA on purpose: /usr/bin/osascript is Apple-signed, so the workflow ships no binary for Gatekeeper to quarantine.
// Syntax stays within ES2020 so older macOS releases' JavaScriptCore can run it.
//
//   timein.js [--format=plain|alfred] <city, landmark or IANA zone>
//
// Test seams: TIMEIN_NOW pins the clock (ISO 8601); TIMEIN_LOOKUPS picks network lookups ("osm,mapkit"; empty = offline);
// TIMEIN_NOMINATIM points OSM lookups at a stand-in server.

ObjC.import('stdlib')

// https://operations.osmfoundation.org/policies/nominatim/ — identify the app, stay under 1 req/s
const NOMINATIM = 'https://nominatim.openstreetmap.org/search?format=json&limit=1&q='
const USER_AGENT = 'alfred-timein (+https://github.com/loginx/alfred-timein)'
const NET_TIMEOUT = 5 // seconds, per network call
const ICON = 'icon.png'
const ERROR_ICON = '/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources/AlertStopIcon.icns'

// Network lookups, most precise first. OSM finds landmarks without location bias; MapKit is the native,
// non-deprecated fallback (biased toward the user's region) for when OSM or CLGeocoder is unavailable.
// A lookup returns a zone name, null when it couldn't answer, or NOT_FOUND to stop the chain: MapKit's
// text search matches almost anything, so asking it after OSM located nothing — or a spot Apple has no
// timezone for, like a summit — would only invent an answer ("South Pole" → America/Chicago).
const LOOKUPS = { __proto__: null, osm, mapkit }
const NOT_FOUND = Symbol('not found')

const FORMATS = {
  __proto__: null,
  plain: { ok: t => t.title, error: fail },
  alfred: {
    ok: t => JSON.stringify({
      items: [{
        title: t.title,
        subtitle: t.subtitle,
        arg: t.title,
        icon: { path: ICON },
        text: { copy: t.title, largetype: t.title },
        action: t.title,
        mods: {
          cmd: { subtitle: 'Copy timezone', arg: t.zone },
          alt: { subtitle: 'Copy ISO 8601 time', arg: t.iso },
        },
        variables: { timezone: t.zone },
      }],
      cache: { seconds: 60 },
      skipknowledge: true,
    }),
    error: msg => JSON.stringify({ items: [{ title: 'Error', subtitle: msg, valid: false, icon: { path: ERROR_ICON } }] }),
  },
}

function run(argv) {
  const args = ObjC.deepUnwrap($.NSProcessInfo.processInfo.arguments)
  const script = args[args.length - argv.length - 1] // osascript [-l JavaScript] <script> ...argv
  const format = argv[0]?.startsWith('--format=') ? argv.shift().slice('--format='.length) : 'plain'
  const out = FORMATS[format] ?? fail(`Unknown format: ${format}`)
  const query = argv.join(' ').trim()
  if (!query) return out.error('City or landmark argument required.')

  const dataDir = env('alfred_workflow_data') ?? `${env('HOME')}/Library/Caches/com.loginx.timein`
  const zone = resolve(query, script.replace(/[^/]*$/, 'capitals.json'), `${dataDir}/cache.json`)
  return zone ? out.ok(clock(zone, new Date(env('TIMEIN_NOW') ?? Date.now()))) : out.error(`Could not geocode: ${query}`)
}

// IANA zone as typed, else user cache, else shipped seed, else network (cached for next time).
// Names this macOS doesn't know (e.g. Europe/Kyiv before tzdata 2022b) count as misses.
function resolve(query, seedPath, cachePath) {
  const lc = query.toLowerCase()
  const typed = ObjC.deepUnwrap($.NSTimeZone.knownTimeZoneNames).find(n => n.toLowerCase() === lc)
  if (typed) return typed

  const cache = readJSON(cachePath)
  const hit = known(cache[lc]) ?? known(readJSON(seedPath)[lc])
  if (hit) return hit

  for (const name of env('TIMEIN_LOOKUPS')?.split(',').filter(Boolean) ?? Object.keys(LOOKUPS)) {
    const found = (LOOKUPS[name] ?? fail(`Unknown lookup: ${name}`))(query)
    if (found === NOT_FOUND) return
    const zone = known(found)
    if (zone) {
      writeJSON(cachePath, { ...cache, [lc]: zone })
      return zone
    }
  }
}

function osm(query) {
  const hits = fetchJSON((env('TIMEIN_NOMINATIM') ?? NOMINATIM) + encodeURIComponent(query))
  if (!Array.isArray(hits)) return null
  if (!hits.length) return NOT_FOUND
  const [hit] = hits
  ObjC.import('CoreLocation')
  // Deprecated in macOS 26, but its replacement (MKReverseGeocodingRequest) never calls back under JXA
  const geocoder = cls('CLGeocoder')
  if (!some(geocoder)) return null
  const location = cls('CLLocation').alloc.initWithLatitudeLongitude(+hit.lat, +hit.lon)
  const [marks] = settle(done => geocoder.alloc.init.reverseGeocodeLocationCompletionHandler(location, done))
  return zoneName(some(marks)?.firstObject.timeZone) || NOT_FOUND
}

function mapkit(query) {
  ObjC.import('MapKit')
  const req = cls('MKLocalSearchRequest').alloc.init
  req.naturalLanguageQuery = query
  const [res] = settle(done => cls('MKLocalSearch').alloc.initWithRequest(req).startWithCompletionHandler(done))
  return zoneName(some(res)?.mapItems.firstObject.timeZone)
}

function fetchJSON(url) {
  const req = cls('NSMutableURLRequest').requestWithURL($.NSURL.URLWithString(url))
  req.setValueForHTTPHeaderField(USER_AGENT, 'User-Agent')
  req.timeoutInterval = NET_TIMEOUT
  const [data, res] = settle(done => cls('NSURLSession').sharedSession.dataTaskWithRequestCompletionHandler(req, done).resume)
  if (!some(data) || Number(res.statusCode) !== 200) return null // JXA returns 64-bit ObjC integers as strings
  try {
    return JSON.parse($.NSString.alloc.initWithDataEncoding(data, $.NSUTF8StringEncoding).js)
  } catch {
    return null
  }
}

// Fixed-format en_US_POSIX output is stable across macOS releases (Apple QA1480); Intl/ICU's is not
// ("Oct 7 at 7:41 PM" vs "Oct 7, 7:41 PM").
function clock(zone, now) {
  const tz = $.NSTimeZone.timeZoneWithName(zone)
  const fmt = $.NSDateFormatter.alloc.init
  const iso = $.NSISO8601DateFormatter.alloc.init
  fmt.locale = $.NSLocale.localeWithLocaleIdentifier('en_US_POSIX')
  fmt.dateFormat = 'EEE, MMM d, h:mm a'
  fmt.timeZone = iso.timeZone = tz
  return {
    zone,
    title: `${zone} - ${fmt.stringFromDate(now).js}`,
    subtitle: `Current time in ${zone.split('/').pop().replace(/_/g, ' ')} (${tz.abbreviationForDate(now).js})`,
    iso: iso.stringFromDate(now).js,
  }
}

// JXA has no await: spin the run loop until the completion handler fires or time runs out.
function settle(start) {
  let args
  start((...a) => { args = a })
  const until = $.NSDate.dateWithTimeIntervalSinceNow(NET_TIMEOUT)
  while (!args && until.timeIntervalSinceNow > 0)
    $.NSRunLoop.currentRunLoop.runModeBeforeDate($.NSDefaultRunLoopMode, $.NSDate.dateWithTimeIntervalSinceNow(0.02))
  return args ?? []
}

// Prototype-free, so a query like "constructor" can't find Object.prototype members; missing or corrupt → empty.
function readJSON(path) {
  const map = Object.create(null)
  try {
    return Object.assign(map, JSON.parse($.NSString.stringWithContentsOfFileEncodingError(path, $.NSUTF8StringEncoding, $()).js))
  } catch {
    return map
  }
}

function writeJSON(path, value) {
  $.NSFileManager.defaultManager.createDirectoryAtPathWithIntermediateDirectoriesAttributesError($(path).stringByDeletingLastPathComponent, true, $(), $())
  $(JSON.stringify(value, null, 2)).writeToFileAtomicallyEncodingError(path, true, $.NSUTF8StringEncoding, $())
}

function fail(msg) {
  console.log(`Error: ${msg}`) // osascript sends console.log to stderr
  $.exit(1)
}

// `$.Name` resolution depends on framework import order ($.NSURLSession is undefined under Foundation alone)
function cls(name) { return $.NSClassFromString(name) }
function env(name) { return ObjC.unwrap($.NSProcessInfo.processInfo.environment.objectForKey(name)) }
// Strings only: NSTimeZone throws on anything else (lookups return null; cache.json may be hand-edited)
function known(zone) { return typeof zone === 'string' && some($.NSTimeZone.timeZoneWithName(zone)) ? zone : null }
function some(obj) { return obj && !obj.isNil() ? obj : null }
function zoneName(tz) { return some(tz) ? tz.name.js : null }
